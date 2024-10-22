package repository

import (
	"context"
	"muxproject1/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Employee interface{
	InsertEmployee(emp *model.Employee) (interface{}, error)
	FindEmployeeByID(emp_id string) (*model.Employee, error)
	FindAllEmployees() ([]model.Employee, error)
	UpdateEmployee(emp_id string, updateEmp model.Employee) (int64, error)
	DeleteEmployeeByID(emp_id string) (int64,error)
	DeleteAllEmployees() (int64,error)
}
type EmployeeRepo struct {
	MongoCollection *mongo.Collection
}

func NewRepo(coll *mongo.Collection) Employee{
	emp:= &EmployeeRepo{
		MongoCollection:coll,
	}
	return emp
}

func (r *EmployeeRepo) InsertEmployee(emp *model.Employee) (interface{}, error) {
	result, err := r.MongoCollection.InsertOne(context.TODO(), emp)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (r *EmployeeRepo) FindEmployeeByID(emp_id string) (*model.Employee, error) {
	var emp model.Employee
	
	filter := bson.D{{Key: "employee_id", Value: emp_id}}
	err := r.MongoCollection.FindOne(context.TODO(), filter).Decode(&emp)
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

func (r *EmployeeRepo) FindAllEmployees() ([]model.Employee, error) {
	cursor, err := r.MongoCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var employees []model.Employee
	err = cursor.All(context.TODO(), &employees)
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func (r *EmployeeRepo) UpdateEmployee(emp_id string, updateEmp model.Employee) (int64, error) {
	
	filter := bson.D{{Key: "employee_id", Value: emp_id}}
	
	update := bson.D{{ Key: "$set", Value:updateEmp}}
	updatResult,err := r.MongoCollection.UpdateOne(context.TODO(), filter, update)

	// update := bson.D{{
	// 	Key: "$set",
	// 	Value:
	// 	 bson.D{
	// 		{Key: "employee_id", Value: updateEmp.EmployeeID}, {Key: "name", Value: updateEmp.Name},
	// 		{Key: "department", Value: updateEmp.Department},
	// 	}}}
	// updateOptions := options.Update().SetUpsert(true)
	// updatResult,err := r.MongoCollection.UpdateOne(context.TODO(), filter, update)
	
	if err != nil {
		return 0, err
	}
	return updatResult.ModifiedCount, nil
}

func (r *EmployeeRepo) DeleteEmployeeByID(emp_id string) (int64,error){
	result,err:= r.MongoCollection.DeleteOne(context.TODO(),bson.D{{Key:"employee_id",Value:emp_id}})
	if(err!=nil){
		return 0,err
	}
	return result.DeletedCount,nil

} 

func (r *EmployeeRepo) DeleteAllEmployees() (int64,error){
	result,err := r.MongoCollection.DeleteMany(context.TODO(),bson.D{{}})
	if(err!=nil){
		return 0,err
	}
	return result.DeletedCount,nil
}

