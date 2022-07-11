package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "8084"
	rpcPort  = "5001"
	mongoUrl = "mongodb://mongo:27017"
	gRpc     = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// const webPort = "8084"

	log.Fatal("working from main in the main file")

	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}
}

// func (app *Config) serve() {
// 	// log.Fatal("working from serve in main file")

// 	srv := &http.Server{
// 		Addr:    fmt.Sprintf(":%s", webPort),
// 		Handler: app.routes(),
// 	}
// 	err := srv.ListenAndServe()
// 	if err != nil {
// 		log.Panic()
// 	}

// }

func connectToMongo() (*mongo.Client, error) {
	// const mongoUrl = "mongodb://mongo:27017"

	// log.Fatal("working from connectToMongo in main file")

	clientOptions := options.Client().ApplyURI(mongoUrl)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("error connecting", err)
		return nil, err
	}

	return c, nil
}
