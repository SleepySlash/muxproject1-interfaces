package signin

import (
	"encoding/json"
	"log"
	"muxproject1/AdminRepository"
	"muxproject1/auth"
	"muxproject1/model/adminModel"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type AdminService struct {
	AdminCollection *mongo.Collection
}

type Response struct{
	Data interface{};
	Error string
}

func (s *AdminService) Login(w http.ResponseWriter, r *http.Request){
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
	adminRepo := AdminRepository.AdminRepo{
		MongoCollection: s.AdminCollection,
	}
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



func (s *AdminService) Register(w http.ResponseWriter, r *http.Request){
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
	adminRepo := AdminRepository.AdminRepo{
		MongoCollection: s.AdminCollection,
	}
	
	res,err := adminRepo.CreateAdmin(admin)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error registering", err)
		result.Error = err.Error()
		return 
	}
	result.Data=res

}