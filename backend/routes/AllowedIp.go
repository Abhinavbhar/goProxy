// AllowedIp handler - returns all IPs from ipCollection
package routes

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func AllowedIp(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := ipCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch IPs"})
	}
	defer cursor.Close(ctx)

	var ips []string
	for cursor.Next(ctx) {
		var entry IPEntry
		if err := cursor.Decode(&entry); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to decode IP"})
		}
		ips = append(ips, entry.IP)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"ips":     ips,
	})
}
