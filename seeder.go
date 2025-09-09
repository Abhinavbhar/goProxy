package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IPEntry model (same as your backend)
type IPEntry struct {
	IP        string    `bson:"ip"`
	CreatedAt time.Time `bson:"created_at"`
}

func main() {
	// Connect to Mongo
	client, err := mongo.Connect(context.TODO(),
		options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Mongo connect error:", err)
	}
	defer client.Disconnect(context.TODO())

	// Use your DB + collection
	db := client.Database("mydb")
	ipCollection := db.Collection("ips")

	// Insert 10,000 fake IPs
	var docs []interface{}
	totalInserted := 0

	for i := 1; i <= 10000; i++ {
		// Generate IPs across multiple subnets to avoid running out
		// This creates IPs like: 192.168.1.1-255, 192.168.2.1-255, etc.
		subnet := ((i - 1) / 254) + 1
		host := ((i - 1) % 254) + 1
		ip := fmt.Sprintf("192.168.%d.%d", subnet, host)

		docs = append(docs, IPEntry{
			IP:        ip,
			CreatedAt: time.Now().Add(time.Duration(-i) * time.Minute), // Vary timestamps
		})

		// Insert in batches of 1000 to avoid memory spike
		if i%1000 == 0 {
			_, err := ipCollection.InsertMany(context.TODO(), docs, options.InsertMany().SetOrdered(false))
			if err != nil {
				fmt.Printf("Batch insert error (continuing): %v\n", err)
			} else {
				totalInserted += len(docs)
			}
			docs = []interface{}{} // reset batch
			fmt.Printf("Inserted batch ending at IP %d (total: %d)\n", i, totalInserted)
		}
	}

	// Insert any remaining documents (in case total isn't divisible by 1000)
	if len(docs) > 0 {
		_, err := ipCollection.InsertMany(context.TODO(), docs, options.InsertMany().SetOrdered(false))
		if err != nil {
			fmt.Printf("Final batch insert error: %v\n", err)
		} else {
			totalInserted += len(docs)
		}
	}

	fmt.Printf("✅ Done! Successfully inserted %d IPs into 'ips' collection\n", totalInserted)
}
