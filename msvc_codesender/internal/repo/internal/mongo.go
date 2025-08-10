package repo_internal

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongo(ctx context.Context, uri string, dbName string) (*Mongo, error) {

	// таймаут 15 сек если ctx не придет быстрее
	ctxTimeout, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctxTimeout, clientOpts)
	if err != nil {
		return nil, err
	}

	// проверяем подключение
	if err := client.Ping(ctxTimeout, nil); err != nil {
		return nil, err
	}

	db := client.Database(dbName)

	return &Mongo{
		Client:   client,
		Database: db,
	}, nil
}

func (m *Mongo) Close(ctx context.Context) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	m.Client.Disconnect(ctxTimeout)
}

func (m *Mongo) LogResponse(ctx context.Context, gatewayName string, response string) {
	collection := m.Database.Collection(gatewayName) // коллекция создастся автоматические если ее нет

	doc := bson.M{
		"response":   response,
		"created_at": time.Now(),
	}

	_, err := collection.InsertOne(ctx, doc)
	if err != nil {
		log.Printf("mongo log insert failed: %v", err)
	}
}
