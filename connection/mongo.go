package connection

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoConnection() (*mongo.Database, error) {
	host := os.Getenv("MONGO_HOSTIP")
	port := os.Getenv("MONGO_PORT")
	user := os.Getenv("MONGO_USER")
	pass := os.Getenv("MONGO_PASSWORD")
	dbName := os.Getenv("MONGO_DB")

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", user, pass, host, port)
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return client.Database(dbName), nil
}
