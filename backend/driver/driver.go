package driver

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// DB holds the connection to the database
type DB struct {
	Conn *mongo.Database
}

var dbConn = &DB{}

func ConnectMongoDB(uri string, dbname string) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	cli, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Ping to check if connection was successful
	if err := cli.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}

	dbConn.Conn = cli.Database(dbname)

	return dbConn, nil
}
