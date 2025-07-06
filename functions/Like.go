package functions

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// LikePublication add a like to a publication ..
func LikePublication(db *mongo.Database, publicationID string, userID int) (string, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	coll := db.Collection("Publications")

	pubID, err := primitive.ObjectIDFromHex(publicationID)
	if err != nil {
		return "", http.StatusBadRequest, errors.New("Invalid publication ID")
	}

	userIDstr := intToString(userID)

	// If the user not in Likes
	filter := bson.M{
		"_id":   pubID,
		"Likes": bson.M{"$ne": userIDstr},
	}
	update := bson.M{
		"$push": bson.M{"Likes": userIDstr},
	}

	res, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", http.StatusInternalServerError, errors.New("Database error: " + err.Error())
	}
	if res.MatchedCount == 0 {
		return "", http.StatusBadRequest, errors.New("Yo have already liked this publication")
	}

	return "Like Added like", http.StatusOK, nil
}

// intToString
func intToString(n int) string {
	return fmt.Sprintf("%d", n)
}
