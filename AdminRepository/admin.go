package AdminRepository

import (
	"context"
	"muxproject1/model/adminModel"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminRepo struct {
	MongoCollection *mongo.Collection

}

// In AdminRepository package
func (r *AdminRepo) CreateAdmin(adm adminModel.Admin) (interface{}, error) {
	id := uuid.New().String()
	adm.AdminID = id
    res, err := r.MongoCollection.InsertOne(context.TODO(), adm)
    if err != nil {
        return nil, err
    }
    return res.InsertedID, err
}

func (r *AdminRepo) GetAdminByPassword(admin adminModel.Admin) (bool, error) {
    filter := bson.D{{Key: "password", Value: admin.Password}}
    err := r.MongoCollection.FindOne(context.TODO(), filter).Decode(&admin)
    if err != nil {
        return false, err
    }
    return true, nil
}
