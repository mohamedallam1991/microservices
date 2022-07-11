package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	log.Fatal("new in models")
	client = mongo
	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	log.Fatal("Insert in models")

	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error inserting into logs")
		return err
	}
	return err
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	log.Fatal("All in models")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := client.Database("logs").Collection("logs")
	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})
	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("error finding all docs error", err)
		return nil, err
	}
	defer cursor.Close(ctx)
	var logs []*LogEntry
	for cursor.Next(ctx) {
		var item LogEntry
		err := cursor.Decode(&item)
		if err != nil {
			log.Println("Error with decoding the cursor", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}
	}
	return logs, nil
}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	log.Fatal("GetOne in models")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := client.Database("logs").Collection("logs")
	DocID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": DocID}).Decode(&entry)

	if err != nil {
		return nil, err
	}
	return &entry, nil
}

// DropCollection deletes the logs collection and everything in it
func (l *LogEntry) DropCollection() error {
	log.Fatal("DropCollection in models")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}

// Update updates on record, by id
func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	log.Fatal("Update in models")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	docID, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return nil, err
	}

	log.Println("Matching", l.ID)
	collection := client.Database("logs").Collection("logs")

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{"$set", bson.D{
				{"name", l.Name},
				{"data", l.Data},
				{"updated_at", time.Now()},
			}},
		},
	)
	log.Println("Matched:", result.MatchedCount)
	log.Println("Modified:", result.ModifiedCount)
	if err != nil {
		return nil, err
	}

	return result, nil
}
