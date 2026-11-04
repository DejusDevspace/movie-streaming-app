package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/DejusDevspace/movie-streaming-app/Server/MagicStreamMoviesServer/database"
	"github.com/DejusDevspace/movie-streaming-app/Server/MagicStreamMoviesServer/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var movieCollection *mongo.Collection = database.Opencollection("movies")

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
