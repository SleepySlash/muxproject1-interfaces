package signin

import (
	"encoding/json"
	"log"
	"muxproject1/AdminRepository"
	"muxproject1/auth"
	"muxproject1/model/adminModel"
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/source/file" // Import the file source driver
	_ "github.com/lib/pq"                                // PostgreSQL driver
)

type Service interface{
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
}
type adminService struct {
	AdminCollection AdminRepository.AdminService
}

type Response struct{
	Data interface{}
	Error string
}

func NewService() Service{
	return &adminService{
		AdminCollection: AdminRepository.NewService(),
	}
}

func (s *adminService) Login(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	log.Println("starting the login")
	result :=  &Response{}
	defer json.NewEncoder(w).Encode(result)
	
	var adm struct{
		Name     string `json:"name,omitempty"`
		Password string `json:"password,omitempty"`
	}
	err := json.NewDecoder(r.Body).Decode(&adm)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid credentials", err)
		result.Error = err.Error()
		return 
	}

	var admin adminModel.Admin
	admin.Name = adm.Name
	admin.Password = adm.Password
	
	okay,err := s.AdminCollection.GetAdminByPassword(admin)
	if err !=nil{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error logging in", err)
		result.Error = err.Error()
		return 
	}
	log.Println("Admin exists, logging in...")
	if okay {
		tokenString, err := auth.CreateToken(admin.Name)
		if err != nil {
		   w.WriteHeader(http.StatusInternalServerError)
		   log.Println("Error creating token")
		   return
		}
	
		w.Header().Set("Authorization", "Bearer "+tokenString)  // Corrected here
		result.Data = admin
		w.WriteHeader(http.StatusOK)
		log.Println("Login successful for ", admin.Name)
	}	
}



func (s *adminService) Register(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	log.Println("let's register")
	result :=  &Response{}
	defer json.NewEncoder(w).Encode(result)
	
	var adm struct{
		Name     string `json:"name,omitempty"`
		Password string `json:"password,omitempty"`
	}
	err := json.NewDecoder(r.Body).Decode(&adm)
	admin := *adminModel.NewAdmin(adm.Name,adm.Password)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid credentials", err)
		result.Error = err.Error()
		return 
	}
	
	err = s.AdminCollection.CreateAdmin(admin)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error registering", err)
		result.Error = err.Error()
		return 
	}
	result.Data=admin
}