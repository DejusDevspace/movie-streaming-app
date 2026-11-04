package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func DBInstance() *mongo.Client {
	// Load env file (local development)
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: unable to find .env file.")
	}

	// Get MongoDB URI from environment variable
	MongoDb := os.Getenv("MONGODB_URI")
	if MongoDb == "" {
		log.Fatal("MONGODB_URI not set!")
	}

	// Set client options and connect to MongoDB
	clientOptions := options.Client().ApplyURI(MongoDb)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil
	}

	return client
}

// Client object to hold the MongoDB client instance
var Client *mongo.Client = DBInstance()

func Opencollection(collectionName string) *mongo.Collection {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: unable to find .env file.")
	}

	// Get database name from environment variable
	databaseName := os.Getenv("DATABASE_NAME")
	if databaseName == "" {
		log.Fatal("DATABASE_NAME not set!")
	}

	collection := Client.Database(databaseName).Collection(collectionName)
	if collection == nil {
		return nil
	}

	return collection
}
