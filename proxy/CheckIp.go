package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ipCollection  *mongo.Collection
	mongoClient   *mongo.Client
	mongoDatabase = "mydb"
	mongoURI      = "mongodb://localhost:27017"
)

func InitMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}

	mongoClient = client
	ipCollection = client.Database(mongoDatabase).Collection("ips")
}

// CheckIp checks if IP is allowed
func CheckIp(ip string) bool {

	_, ok := ipBandwidth[ip]
	if ok {
		return true
	}

	// Not in memory → query MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	count, err := ipCollection.CountDocuments(ctx, bson.M{"ip": ip})
	if err != nil {
		fmt.Println("Mongo query error:", err)
		return false
	}

	if count > 0 {
		// Add to memory
		ipMutex.Lock()
		ipBandwidth[ip] = 0
		ipMutex.Unlock()
		return true
	}

	return false
}
