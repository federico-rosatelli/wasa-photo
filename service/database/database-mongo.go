package database

import (
	"context"
	"wasa-photo/service/api/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppDatabaseMongo interface {
	GetName() (string, error)
	SetName(name string) error

	UpdateOne(collection mongo.Collection, filter primitive.D, query primitive.D) error
	UpdateOnePush(collection mongo.Collection, filter primitive.D, query primitive.M) error
	InsertOne(collection mongo.Collection, query primitive.M) error
	GetUsersCollection() mongo.Collection
	GetProfilesCollection() mongo.Collection
	GetSessionsCollection() mongo.Collection
	// FindOne(collection mongo.Collection, query string) error

	Ping() error
}

type appdbmongo struct {
	c        *mongo.Client
	users    *mongo.Collection
	profiles *mongo.Collection
	sessions *mongo.Collection
}

// var CollectionUsers *mongo.Collection
// var CollectionProfiles *mongo.Collection
// var CollectionSessions *mongo.Collection
var Ctx = context.TODO()

func MakeInit(client *mongo.Client) (AppDatabaseMongo, error) {
	collectionUsers := client.Database("wasa").Collection("users")
	if collectionUsers == nil {
		return nil, errors.NewErrStatus("Error Creating users Collection")
	}
	collectionProfiles := client.Database("wasa").Collection("profiles")
	if collectionProfiles == nil {
		return nil, errors.NewErrStatus("Error Creating profiles Collection")
	}
	collectionSessions := client.Database("wasa").Collection("sessions")
	if collectionSessions == nil {
		return nil, errors.NewErrStatus("Error Creating sessions Collection")
	}
	return &appdbmongo{
		c:        client,
		users:    collectionUsers,
		profiles: collectionProfiles,
		sessions: collectionSessions,
	}, nil
}

func (db *appdbmongo) Ping() error {
	return db.c.Ping(Ctx, nil)
}

func (db *appdbmongo) GetName() (string, error) {
	var name struct {
		Username string
	}
	err := db.users.FindOne(Ctx, bson.D{{}}).Decode(&name)
	return name.Username, err
}

func (db *appdbmongo) SetName(name string) error {
	var str struct {
		Username string
	}
	str.Username = name
	_, err := db.users.InsertOne(Ctx, str)
	return err
}

func (db *appdbmongo) GetUsersCollection() mongo.Collection {
	return *db.users
}

func (db *appdbmongo) GetProfilesCollection() mongo.Collection {
	return *db.profiles
}

func (db *appdbmongo) GetSessionsCollection() mongo.Collection {
	return *db.sessions
}

// func NewUser(user User) error {
// 	_, err := collectionUsers.InsertOne(ctx, user)
// 	return err
// }
