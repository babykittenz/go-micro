package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var client *mongo.Client

// New initializes a Models instance using the provided MongoDB client and sets the global client variable.
func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}

// Models is a type that aggregates various data structures and their associated methods for managing application models.
type Models struct {
	LogEntry LogEntry
}

// LogEntry represents a log entry stored in the database, including metadata such as creation and update timestamps.
type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreateAt  time.Time `bson:"create_at" json:"create_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// Insert adds a new LogEntry into the database and sets creation and update timestamps. Returns an error if insertion fails.
func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreateAt:  time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// All retrieves all LogEntry documents from the database, sorted by creation date in descending order. Returns a slice of LogEntry pointers and an error.
func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{{
		"created_at", -1,
	}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry
	for cursor.Next(ctx) {
		var item LogEntry
		if err != nil {
			log.Println("error decoding log into slice", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}

	}
	return logs, nil
}

// GetOne retrieves a single LogEntry from the database by its unique identifier. Returns the LogEntry and an error if any.
func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var entry LogEntry

	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &entry, nil
}

// DropCollection removes the "logs" collection from the database. Returns an error if the operation fails.
func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	if err := collection.Drop(ctx); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Update modifies an existing LogEntry document in the database by its ID and updates its fields. Returns the update result and an error.
func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{{"$set", bson.D{
			{"name", l.Name},
			{"data", l.Data},
			{"updated_at", time.Now()},
		},
		}},
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result, nil
}
