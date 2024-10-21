package repository

import (
	"context"
	"log"
	"muxproject1/model"
	"testing"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func newMongoClient() *mongo.Client{

	clientOptions := options.Client().ApplyURI("mongodb+srv://hiimindra12345:a0gIMLwmnCd3R00a@cluster0.93xdm.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")
	mongoTestClient,err := mongo.Connect(context.TODO(),clientOptions)
	
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	err = mongoTestClient.Ping(context.TODO(),readpref.Primary())
	if err != nil {
		log.Fatal("ping error ",err)
	}
	log.Println("Ping succesfull!")
	return mongoTestClient
}

func TestMongoOperations(t *testing.T){
 
	mongoTestClient := newMongoClient()

	defer mongoTestClient.Disconnect(context.TODO())

	//dummy data
	emp1 := uuid.New().String()
	emp2 := uuid.New().String()

	coll := mongoTestClient.Database("companyDB").Collection("employee_test")

	empRepo := EmployeeRepo{MongoCollection: coll}


	// insert employee 1
	
	t.Run("insert empolyee 2", func(t *testing.T){
		emp:=model.Employee{
			Name: "Steve Rogers",
			Department: "HR",
			EmployeeID: emp2,
		}
		result,err := empRepo.InsertEmployee(&emp)
		if(err!=nil){
			t.Fatal(err)
		}
		log.Println("Inserted Employee 2 with ID: ",result)
	})
	
	t.Run("insert employee 1", func(t *testing.T){
		emp:= model.Employee{
			Name: "Tony Stark",
			Department: "Engineering",
			EmployeeID: emp1,
		}
		result,err := empRepo.InsertEmployee(&emp)

		if(err!=nil){
			t.Fatal(err)
		}
		log.Println("Inserted Employee 1 with ID: ",result)
	})
	t.Run("get emp 1", func(t *testing.T){
		result,err := empRepo.FindEmployeeByID(emp1)
		if(err!=nil){
			t.Fatal(err)
		}
		log.Println("emp 1: ",result.Name)
	})
	
	t.Run("get all employees", func(t *testing.T){
		result,err := empRepo.FindAllEmployees()
		if(err!=nil){
			t.Fatal(err)
		}
		log.Println("emp 1: ",result)
	})
	
	t.Run("Update employee 1 Name",func(t *testing.T){
		emp := model.Employee{
			Name: "Iron Man",
			Department: "Super Hero",
			EmployeeID: "8cb02b92-4f5c-4ef7-b302-d72df4a327d0",
		}
		
		result,err := empRepo.UpdateEmployee("8cb02b92-4f5c-4ef7-b302-d72df4a327d0",emp)
		if err != nil {
			t.Fatal(err)
		}
		log.Println("Updated employee 1: ",result)
	})
	
	t.Run("Delete employee 1",func(t *testing.T){
			result,err := empRepo.DeleteEmployeeByID(emp1)
			if err != nil {
			t.Fatal(err)
		}
		log.Println("Deleted employee 1: ",result)
	})
}