package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/DejusDevspace/movie-streaming-app/Server/MagicStreamMoviesServer/database"
	"github.com/DejusDevspace/movie-streaming-app/Server/MagicStreamMoviesServer/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Database connection and collection initialization
var movieCollection *mongo.Collection = database.Opencollection("movies")

// Validator instance for input validation
var validate = validator.New()

func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movies []models.Movie

		// Fetch all movies from the collection
		cursor, err := movieCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies."})
		}
		// Close the cursor after use (memory management)
		defer cursor.Close(ctx)

		// Decode the cursor into the movies slice
		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movies."})
		}

		// Pass the movies slice (JSON response) to the client
		c.JSON(http.StatusOK, movies)
	}
}

func GetMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Implementation for fetching a single movie
		movieID := c.Param("imdb_id")
		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
			return
		}

		var movie models.Movie

		// Fetch the movie from the collection based on the provided IMDb ID
		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		c.JSON(http.StatusOK, movie)
	}
}

func AddMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movie models.Movie

		// Bind the JSON payload to the movie struct
		if err := c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Validate the movie struct using the validator package
		if err := validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}

		// Insert the data into the db
		result, err := movieCollection.InsertOne(ctx, movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add movie"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Movie added successfully", "id": result.InsertedID})
	}
}
