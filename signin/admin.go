package signin

import (
	"encoding/json"
	"log"
	"muxproject1/AdminRepository"
	"muxproject1/auth"
	"muxproject1/model/adminModel"
	"net/http"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service interface{
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
}
type adminService struct {
	AdminCollection *mongo.Collection
}

type Response struct{
	Data interface{}
	Error string
}

func NewService(collection *mongo.Collection) Service{
	svc:= &adminService{
		AdminCollection: collection,
	}
	return svc
}

func (s *adminService) Login(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	log.Println("starting the login")
	result :=  &Response{}
	defer json.NewEncoder(w).Encode(result)
	
	var admin adminModel.Admin
	err := json.NewDecoder(r.Body).Decode(&admin)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid credentials", err)
		result.Error = err.Error()
		return 
	}
	adminRepo := AdminRepository.NewService(s.AdminCollection)
	okay,err := adminRepo.GetAdminByPassword(admin)
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
	
	var admin adminModel.Admin
	err := json.NewDecoder(r.Body).Decode(&admin)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid credentials", err)
		result.Error = err.Error()
		return 
	}
	adminRepo := AdminRepository.NewService(s.AdminCollection)
	id := uuid.NewString()
	admin.AdminID = id
	res,err := adminRepo.CreateAdmin(admin)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error registering", err)
		result.Error = err.Error()
		return 
	}
	result.Data=res

}