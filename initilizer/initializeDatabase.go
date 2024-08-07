package initilizer

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitializeDatabase(client *mongo.Client, ctx context.Context) *mongo.Database {
	db := client.Database("crud")

	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"email", "password", "role"},
			"title":    "User Object Validation",
			"properties": bson.M{
				"email": bson.M{
					"bsonType":    "string",
					"description": "must be a string and is required",
				},
				"password": bson.M{
					"bsonType":    "string",
					"description": "must be a string and is required",
				},
				"role": bson.M{
					"enum":        []string{"user", "admin"},
					"description": "must be user or admin and is required",
				},
				"age": bson.M{
					"bsonType":    "int",
					"minimum":     0,
					"description": "must be an integer greater than or equal to 0",
				},
			},
		},
	}

	err := db.CreateCollection(ctx, "users", options.CreateCollection().SetValidator(validator))
	_, ok := err.(mongo.CommandError)
	if err != nil && !ok {
		log.Fatal("Failed to create collection in mongodb, error: " + err.Error())
	}

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = db.Collection("users").Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Fatal("Failed to create index in mongodb,error: " + err.Error())
	}
	return db
}
