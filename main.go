package main

import (
	"context"
	"log"
	"muxproject1/signin"
	"muxproject1/usecase"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoClient *mongo.Client

func init(){
	var err error
	err=godotenv.Load()
	if err!=nil{
		log.Fatal("error at the dot env load",err)
	}
	log.Println("Loaded env variables")
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	mongoClient,err = mongo.Connect(context.Background(),clientOptions)
	if err!=nil{
		log.Fatal("error at connecting",err)
	}
	err = mongoClient.Ping(context.TODO(),readpref.Primary())
	if err!=nil{
		log.Fatal("ping failed",err)
	}
	log.Println("Connected to MONGODB")
}

func main() {
	defer mongoClient.Disconnect(context.Background())

	// Service Variable for Employee
	svc := usecase.NewService(mongoClient.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("COLLECTION_NAME")))
	adm := signin.NewService(mongoClient.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("ADMIN_COLLECTION")))
	
	// Service Variable for Admin
	r := mux.NewRouter()
	r.HandleFunc("/health",healthHandler).Methods("GET")


	r.HandleFunc("/login",adm.Login).Methods("POST")
	r.HandleFunc("/register",adm.Register).Methods("POST")


	r.HandleFunc("/employee",svc.CreateEmployee).Methods("POST")
	r.HandleFunc("/employee/{id}",svc.GetEmployeeByID).Methods("GET")
	r.HandleFunc("/employee",svc.GetAllEmployees).Methods("GET",)
	r.HandleFunc("/employee/{id}",svc.UpdateEmployeeByID).Methods("PUT")
	r.HandleFunc("/employee/delete/{id}", svc.DeleteEmployeeByID).Methods("DELETE")
	r.HandleFunc("/employee/delete",svc.DeleteAllEmployees).Methods("DELETE")

	log.Println("server is on 3000")
	http.ListenAndServe(":3000",r)
}

func healthHandler(w http.ResponseWriter,r *http.Request){
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("runnig..."))	
}