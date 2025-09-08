package routes

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func VerifyAuth(c *fiber.Ctx) error {
	fmt.Println("verifyAuth")
	type ReqBody struct {
		Token string `json:"token"`
	}

	var body ReqBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	// 1. Decode JWT
	secret := "ifYouSeeItuRDead"
	tok, err := jwt.Parse(body.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !tok.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "error": "invalid token"})
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || claims["email"] == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "error": "invalid token claims"})
	}
	email := claims["email"].(string)

	// 2. Get client IP
	clientIP := c.IP()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 3. Find user
	var user User
	err = userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "error": "user not found"})
	}

	// 4. Check if IP exists in user IPs array
	ipExists := false
	for _, ip := range user.IPs {
		if ip == clientIP {
			ipExists = true
			break
		}
	}

	// If IP not present, add it to user document
	update := bson.M{}
	if !ipExists {
		update["$push"] = bson.M{"ips": clientIP}
	}
	// Always update last login
	update["$set"] = bson.M{"last_login": time.Now()}

	if len(update) > 0 {
		_, _ = userCollection.UpdateOne(ctx, bson.M{"email": email}, update)
	}

	// 5. Add IP to IPs collection for fast lookup
	_, _ = ipCollection.UpdateOne(ctx,
		bson.M{"ip": clientIP},
		bson.M{"$setOnInsert": IPEntry{
			IP:        clientIP,
			CreatedAt: time.Now(),
		}},
		options.Update().SetUpsert(true),
	)

	// 6. Respond with user info
	return c.JSON(fiber.Map{
		"success": true,
		"email":   email,
		"ips":     append(user.IPs, clientIP), // return updated list of IPs
	})
}
