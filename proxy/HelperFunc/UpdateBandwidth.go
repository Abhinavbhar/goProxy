package helperfunc

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// Your existing models

func UpdateBandwidth() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Replace with your Users collection name

	var updatedCount int
	var userUpdatedCount int

	// Iterate through your in-memory bandwidth map
	for ip, bandwidth := range IpBandwidth {
		// Update bandwidth in IPEntry collection
		filter := bson.M{"ip": ip}
		update := bson.M{"$set": bson.M{"bandwidth": bandwidth}}

		result, err := ipCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Printf("Error updating IP %s: %v", ip, err)
			continue
		}

		if result.ModifiedCount > 0 {
			updatedCount++
		}

		// Find and update all users with this IP
		userFilter := bson.M{"ips": ip}
		userUpdate := bson.M{"$set": bson.M{"bandwidth": bandwidth}}

		userResult, err := userCollection.UpdateMany(ctx, userFilter, userUpdate)
		if err != nil {
			log.Printf("Error updating users with IP %s: %v", ip, err)
			continue
		}

		userUpdatedCount += int(userResult.ModifiedCount)
	}

	log.Printf("Updated bandwidth for %d IP entries and %d users", updatedCount, userUpdatedCount)
	return nil
}

// Function to set up hourly updates using a ticker
func StartBandwidthUpdater() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	// Run immediately on start
	if err := UpdateBandwidth(); err != nil {
		log.Printf("Error in initial bandwidth update: %v", err)
	}

	// Then run every hour
	for range ticker.C {
		if err := UpdateBandwidth(); err != nil {
			log.Printf("Error in scheduled bandwidth update: %v", err)
		}
	}
}
