package AdminRepository

import (
	"context"
	"muxproject1/model/adminModel"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminService interface{
    CreateAdmin(adm adminModel.Admin) (interface{}, error)
    GetAdminByPassword(admin adminModel.Admin) (bool, error)
}
type adminRepo struct {
	MongoCollection *mongo.Collection
}

func NewService(collection *mongo.Collection) AdminService{
    repo := &adminRepo{
        MongoCollection: collection,
    }
    return repo
}
// In AdminRepository package
func (r *adminRepo) CreateAdmin(adm adminModel.Admin) (interface{}, error) {
	
    res, err := r.MongoCollection.InsertOne(context.TODO(), adm)
    if err != nil {
        return nil, err
    }
    return res.InsertedID, err
}

func (r *adminRepo) GetAdminByPassword(admin adminModel.Admin) (bool, error) {
    filter := bson.D{{Key: "password", Value: admin.Password}}
    err := r.MongoCollection.FindOne(context.TODO(), filter).Decode(&admin)
    if err != nil {
        return false, err
    }
    return true, nil
}