package database

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (db *appdbmongo) UpdateOne(collection mongo.Collection, filter primitive.D, query primitive.D) error {
	_, err := collection.UpdateOne(Ctx, filter, query)
	return err
}

func (db *appdbmongo) UpdateOnePush(collection mongo.Collection, filter primitive.D, query primitive.M) error {
	_, err := collection.UpdateOne(Ctx, filter, query)
	return err
}

func (db *appdbmongo) InsertOne(collection mongo.Collection, query primitive.M) error {
	_, err := db.users.InsertOne(Ctx, query)
	return err
}
