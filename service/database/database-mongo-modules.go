package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (db *appdbmongo) UpdateOne(typeCli int, filter primitive.D, query primitive.D) error {
	var collection mongo.Collection
	switch typeCli {
	case 0:
		collection = *db.profiles
	case 1:
		collection = *db.sessions
	case 2:
		collection = *db.users
	}
	_, err := collection.UpdateOne(Ctx, filter, query)
	return err
}

func (db *appdbmongo) UpdateOnePush(typeCli int, filter primitive.D, query primitive.M) error {
	var collection mongo.Collection
	switch typeCli {
	case 0:
		collection = *db.profiles
	case 1:
		collection = *db.sessions
	case 2:
		collection = *db.users
	}
	_, err := collection.UpdateOne(Ctx, filter, query)
	return err
}

func (db *appdbmongo) UpdateOnePushM(typeCli int, filter primitive.M, query primitive.M) error {
	var collection mongo.Collection
	switch typeCli {
	case 0:
		collection = *db.profiles
	case 1:
		collection = *db.sessions
	case 2:
		collection = *db.users
	}
	_, err := collection.UpdateOne(Ctx, filter, query)
	return err
}

func (db *appdbmongo) InsertOne(typeCli int, query primitive.M) error {
	var collection mongo.Collection
	switch typeCli {
	case 0:
		collection = *db.profiles
	case 1:
		collection = *db.sessions
	case 2:
		collection = *db.users
	}
	_, err := collection.InsertOne(Ctx, query)
	return err
}

func (db *appdbmongo) InsertOneProfile(profileStruct Profile) error {

	collection := *db.profiles
	_, err := collection.InsertOne(Ctx, profileStruct)
	return err
}

func (db *appdbmongo) InsertOneUsers(usersStruct User) error {

	collection := *db.users
	_, err := collection.InsertOne(Ctx, usersStruct)
	return err
}

func (db *appdbmongo) InsertOneSession(sessionStruct Session) error {

	collection := *db.sessions
	_, err := collection.InsertOne(Ctx, sessionStruct)
	return err
}

func (db *appdbmongo) BackUpProfiles() (mongo.Cursor, error) {
	collection := *db.profiles
	data, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return *data, err
	}
	return *data, nil
}

func (db *appdbmongo) BackUpUsers() (mongo.Cursor, error) {
	collection := *db.users
	data, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return *data, err
	}
	return *data, nil
}

func (db *appdbmongo) BackUpSessions() (mongo.Cursor, error) {
	collection := *db.sessions
	data, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return *data, err
	}
	return *data, nil
}
