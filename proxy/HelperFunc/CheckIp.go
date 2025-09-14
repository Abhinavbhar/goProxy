package helperfunc

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ipCollection   *mongo.Collection
	userCollection *mongo.Collection
	mongoDatabase  = "mydb"
)

func InitMongo() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		panic("MONGO_URI not set in environment")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Println("failed to connect to mongo")
	}

	fmt.Println("MongoDB connected!")

	ipCollection = client.Database(mongoDatabase).Collection("ips")
	userCollection = client.Database(mongoDatabase).Collection("users")
}

// CheckIp checks if IP is allowed
func CheckIp(ip string) bool {

	_, ok := IpBandwidth[ip]
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
		IpMutex.Lock()
		IpBandwidth[ip] = 0
		IpMutex.Unlock()
		return true
	}

	return false
}
