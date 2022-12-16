package database

import (
	"context"
	customError "wasa-photo/service/api/customErrors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppDatabaseMongo interface {
	// GetName() (string, error)
	// SetName(name string) error

	UpdateOne(typeCli int, filter primitive.D, query primitive.D) error
	UpdateOnePush(typeCli int, filter primitive.D, query primitive.M) error
	UpdateOnePushM(typeCli int, filter primitive.M, query primitive.M) error
	InsertOne(typeCli int, query primitive.M) error
	InsertOneProfile(profileStruct Profile) error
	InsertOneUsers(usersStruct User) error
	InsertOneSession(sessionStruct Session) error
	BackUpProfiles() (mongo.Cursor, error)
	BackUpUsers() (mongo.Cursor, error)
	BackUpSessions() (mongo.Cursor, error)
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
		return nil, customError.NewErrStatus("Error Creating users Collection")
	}
	collectionProfiles := client.Database("wasa").Collection("profiles")
	if collectionProfiles == nil {
		return nil, customError.NewErrStatus("Error Creating profiles Collection")
	}
	collectionSessions := client.Database("wasa").Collection("sessions")
	if collectionSessions == nil {
		return nil, customError.NewErrStatus("Error Creating sessions Collection")
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
