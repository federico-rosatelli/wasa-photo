package database

import (
	"context"
	"web/errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var CollectionUsers *mongo.Collection
var CollectionProfiles *mongo.Collection
var CollectionSessions *mongo.Collection
var Ctx = context.TODO()

func MakeInit() error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(Ctx, clientOptions)
	if err != nil {
		return errors.NewErrStatus("Error Connecting to mongoDB:\n\t" + err.Error())
	}
	CollectionUsers = client.Database("wasa").Collection("users")
	CollectionProfiles = client.Database("wasa").Collection("profiles")
	CollectionSessions = client.Database("wasa").Collection("sessions")
	return nil
}

// func NewUser(user User) error {
// 	_, err := collectionUsers.InsertOne(ctx, user)
// 	return err
// }
