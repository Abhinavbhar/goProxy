package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User model for MongoDB
// User model
type User struct {
	Email     string    `bson:"email"`
	CreatedAt time.Time `bson:"created_at"`
	Bandwidth int       `bson:"bandwidth"`
	LastLogin time.Time `bson:"last_login"`
	IPs       []string  `bson:"ips"` // store all user IPs
}

// IP model for fast lookup (no userEmail)
type IPEntry struct {
	IP        string    `bson:"ip"`         // single IP
	CreatedAt time.Time `bson:"created_at"` // optional for tracking
}

// Google token info response
type TokenInfo struct {
	Email string `json:"email"`
	Aud   string `json:"aud"`
	Exp   string `json:"exp"`
}

// Mongo client (global for now)
var userCollection *mongo.Collection
var ipCollection *mongo.Collection

func InitMongo() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	db := client.Database("mydb")
	userCollection = db.Collection("users")
	ipCollection = db.Collection("ips")

	// Create unique index on IP for fast lookup
	_, _ = ipCollection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.M{"ip": 1},
		Options: options.Index().SetUnique(true),
	})
}

// Login handler
func Login(c *fiber.Ctx) error {
	fmt.Println("Login came")
	type ReqBody struct {
		GoogleToken string `json:"google_token"`
	}
	var body ReqBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	// 1. Verify Google token
	resp, err := http.Get("https://www.googleapis.com/oauth2/v3/tokeninfo?access_token=" + body.GoogleToken)
	if err != nil || resp.StatusCode != 200 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "error": "invalid google token"})
	}
	defer resp.Body.Close()

	var tokenInfo TokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to parse google response"})
	}

	email := tokenInfo.Email
	if email == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "error": "no email found"})
	}

	// Get client IP
	clientIP := c.IP()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existing User
	err = userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&existing)

	if err == mongo.ErrNoDocuments {
		// New user → insert
		newUser := User{
			Email:     email,
			CreatedAt: time.Now(),
			Bandwidth: 0,
			LastLogin: time.Now(),
			IPs:       []string{clientIP},
		}
		_, err = userCollection.InsertOne(ctx, newUser)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to create user"})
		}
	} else {
		// Existing user → update last login & add IP if not exists
		ipExists := false
		for _, ip := range existing.IPs {
			if ip == clientIP {
				ipExists = true
				break
			}
		}
		update := bson.M{"$set": bson.M{"last_login": time.Now()}}
		if !ipExists {
			update["$push"] = bson.M{"ips": clientIP}
		}
		_, _ = userCollection.UpdateOne(ctx, bson.M{"email": email}, update)
	}

	// Insert IP into IP collection for fast lookup (upsert)
	_, _ = ipCollection.UpdateOne(ctx,
		bson.M{"ip": clientIP},
		bson.M{"$setOnInsert": IPEntry{
			IP:        clientIP,
			CreatedAt: time.Now(),
		}},
		options.Update().SetUpsert(true),
	)

	// 3. Create backend JWT (30 days)
	secret := "ifYouSeeItuRDead"
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(30 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to sign token"})
	}

	// 4. Respond
	return c.JSON(fiber.Map{
		"success": true,
		"email":   email,
		"token":   signedToken,
	})
}
