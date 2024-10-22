package usecase

import (
	"encoding/json"
	"fmt"
	"log"
	"muxproject1/auth"
	"muxproject1/model"
	"muxproject1/repository"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)
type Service interface{
	CreateEmployee (w http.ResponseWriter, r *http.Request)
	GetEmployeeByID (w http.ResponseWriter, r *http.Request)
	GetAllEmployees (w http.ResponseWriter, r *http.Request) 
	UpdateEmployeeByID (w http.ResponseWriter, r *http.Request)
	DeleteEmployeeByID (w http.ResponseWriter, r *http.Request)
	DeleteAllEmployees (w http.ResponseWriter, r *http.Request)
}
type EmployeeService struct {
	EmployeeCollection *mongo.Collection
}

type Response struct{
	Data  interface{} `json:"data,omitempty"` 
	Error string      `json:"error,omitempty"`
}

func NewService(client *mongo.Collection) Service{
	svc:=&EmployeeService{
		EmployeeCollection: client,
	}
	return svc
}

func (svc *EmployeeService) CreateEmployee (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	
	tokenString := r.Header.Get("Authorization")
  	if tokenString == "" {
  	  w.WriteHeader(http.StatusUnauthorized)
  	  log.Println("Unauthorized User")
  	  return
  	}
  	tokenString = tokenString[len("Bearer "):]
  	err := auth.VerifyToken(tokenString)
  	if err != nil {
  	  w.WriteHeader(http.StatusUnauthorized)
  	  fmt.Fprint(w, "Invalid token")
  	  return
  	}

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)
	
	var emp model.Employee
	err = json.NewDecoder(r.Body).Decode(&emp)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid body", err)
		res.Error = err.Error()
		return 
	}

	emp.EmployeeID = uuid.New().String()

	empRepo := repository.NewRepo(svc.EmployeeCollection)

	empId,err := empRepo.InsertEmployee(&emp)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error while creating employee",err)
		res.Error= err.Error()
		return
	}
	res.Data = emp.EmployeeID
	w.WriteHeader(http.StatusOK)
	log.Println("employee is inserted with id ",empId,emp)
}

func (svc *EmployeeService) GetEmployeeByID (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
	  w.WriteHeader(http.StatusUnauthorized)
	  log.Println("Unauthorized User")
	  return
	}
	tokenString = tokenString[len("Bearer "):]
	err := auth.VerifyToken(tokenString)
	if err != nil {
	  w.WriteHeader(http.StatusUnauthorized)
	  fmt.Fprint(w, "Invalid token")
	  return
	}
	res := &Response{}
	
	defer json.NewEncoder(w).Encode(res)
	id := mux.Vars(r)["id"]
	empRepo := repository.NewRepo(svc.EmployeeCollection)
	var emp *model.Employee
	emp,err = empRepo.FindEmployeeByID(id)
	if(err!=nil){
		w.WriteHeader(http.StatusNotFound)
		log.Println("error in finding the employee with id",err)
		res.Error = err.Error()
		return
	}
	log.Println("found the employee with id",id,emp)
	res.Data=emp
	w.WriteHeader(http.StatusOK)
}

func (svc *EmployeeService) GetAllEmployees (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
	  w.WriteHeader(http.StatusUnauthorized)
	  log.Println("Unauthorized User")
	  return
	}
	tokenString = tokenString[len("Bearer "):]
	err := auth.VerifyToken(tokenString)
	if err != nil {
	  w.WriteHeader(http.StatusUnauthorized)
	  fmt.Fprint(w, "Invalid token")
	  return
	}
	res := &Response{}
	
	defer json.NewEncoder(w).Encode(res)

	empRepo := repository.NewRepo(svc.EmployeeCollection)

	users,err:= empRepo.FindAllEmployees()
	if err!=nil{
		w.WriteHeader(http.StatusNotFound)
		log.Println("error in finding the employee with id",err)
		res.Error = err.Error()
		return
	}
	log.Println("no of employees in the database are ",len(users))
	log.Println("All the employees in the database are ",users)
	res.Data = users
	w.WriteHeader(http.StatusOK)
}

func (svc *EmployeeService) UpdateEmployeeByID (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
	  w.WriteHeader(http.StatusUnauthorized)
	  log.Println("Unauthorized User")
	  return
	}
	tokenString = tokenString[len("Bearer "):]
	err := auth.VerifyToken(tokenString)
	if err != nil {
	  w.WriteHeader(http.StatusUnauthorized)
	  fmt.Fprint(w, "Invalid token")
	  return
	}
	res := &Response{}
	defer json.NewEncoder(w).Encode(res)
	id := mux.Vars(r)["id"]

	if id==""{
		w.WriteHeader(http.StatusNotFound)
		log.Println("bad user id")
		res.Error ="invalid user id"
		return
	}
	var user model.Employee
	err=json.NewDecoder(r.Body).Decode(&user)
	if err!=nil{
		w.WriteHeader(http.StatusNotFound)
		log.Println("bad user info",err)
		res.Error = err.Error()
		return
	}
	user.EmployeeID = id
	
	empRepo := repository.NewRepo(svc.EmployeeCollection)

	count,err:= empRepo.UpdateEmployee(id,user)

	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("could not update employee",err)
		res.Error = err.Error()
		return
	}
	log.Println("updated employee ",user)
	res.Data=count
	w.WriteHeader(http.StatusOK)
}

func (svc *EmployeeService) DeleteEmployeeByID (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
	  w.WriteHeader(http.StatusUnauthorized)
	  log.Println("Unauthorized User")
	  return
	}
	tokenString = tokenString[len("Bearer "):]
	err := auth.VerifyToken(tokenString)
	if err != nil {
	  w.WriteHeader(http.StatusUnauthorized)
	  fmt.Fprint(w, "Invalid token")
	  return
	}
	res := &Response{}
	defer json.NewEncoder(w).Encode(res)
	id :=mux.Vars(r)["id"]

	if id==""{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("bad user id")
		res.Error ="invalid user id"
		return
	}
	empRepo := repository.NewRepo(svc.EmployeeCollection)
	count,err := empRepo.DeleteEmployeeByID(id)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("could not delete the employee with id",id)
		res.Error = err.Error()
		return
	}
	res.Data = count
	log.Println("delted the user with id",id)
	w.WriteHeader(http.StatusOK)
}


func (svc *EmployeeService) DeleteAllEmployees (w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type","application/json")
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
	  w.WriteHeader(http.StatusUnauthorized)
	  log.Println("Unauthorized User")
	  return
	}
	tokenString = tokenString[len("Bearer "):]
	err := auth.VerifyToken(tokenString)
	if err != nil {
	  w.WriteHeader(http.StatusUnauthorized)
	  fmt.Fprint(w, "Invalid token")
	  return
	}
	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	empRepo := repository.NewRepo(svc.EmployeeCollection)
	count,err := empRepo.DeleteAllEmployees()
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("could not delete the users")
		res.Error = err.Error()
		return
	}
	res.Data = count
	log.Println("total no of employeed deleted is",count)
	w.WriteHeader(http.StatusOK)
}