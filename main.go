package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/NdoleStudio/lemonsqueezy-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"main.go/database"
	"main.go/models"
)

type WebhookPayload struct {
	Meta struct {
		CustomData struct {
			UserID string `json:"user_id"`
		} `json:"custom_data"`
	} `json:"meta"`
	Data struct {
		Attributes struct {
			Status string `json:"status"`
		} `json:"attributes"`
	} `json:"data"`
}

func main() {

	database.ConnectDb()

	client := lemonsqueezy.New(lemonsqueezy.WithAPIKey(os.Getenv("LEMONSQUEEZY_API_KEY")))
	fmt.Println(client)

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})

	fmt.Println("Jai Shree Ram !!")

	imageDir := "./public/images/"
	if _, err := os.Stat(imageDir); os.IsNotExist(err) {
		os.MkdirAll(imageDir, 0755) // Create the directory if it doesn't exist
	}

	app.Post("/webhooks", func(c *fiber.Ctx) error {

		// Retrieve the event name and signature from the headers
		eventType := c.Get("X-Event-Name", "")
		signatureHeader := c.Get("X-Signature", "")

		if eventType == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Missing event type"})
		}

		// Read the request body
		bodyBytes := c.Body()
		// if err != nil {
		// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to read request body"})
		// }

		// Verify the signature
		secret := os.Getenv("LEMONSQUEEZY_WEBHOOKS_SECRET")
		hmacHasher := hmac.New(sha256.New, []byte(secret))
		hmacHasher.Write(bodyBytes)
		expectedSignature := hex.EncodeToString(hmacHasher.Sum(nil))

		if signatureHeader == "" || !hmac.Equal([]byte(expectedSignature), []byte(signatureHeader)) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid signature"})
		}

		// Parse the JSON body
		var payload WebhookPayload
		err := json.Unmarshal(bodyBytes, &payload)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid JSON"})
		}

		// Handle different events based on the event type
		switch eventType {
		case "order_created":
			// Extract user ID and status
			userID := payload.Meta.CustomData.UserID
			isSuccessful := payload.Data.Attributes.Status == "paid"
			log.Printf("Order created: User ID: %s, Successful: %t", userID, isSuccessful)

		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": fmt.Sprintf("Invalid event type: %s", eventType)})
		}

		// If everything is successful, return a positive response
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Webhook received"})
	})

	app.Post("/new", func(c *fiber.Ctx) error {

		file, err := c.FormFile("image")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "No file uploaded",
			})
		}
		// Generate a unique filename and save the file to the image directory
		fileName := filepath.Join(imageDir, fmt.Sprintf("%s", file.Filename))
		if err := c.SaveFile(file, fileName); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "File saving failed",
			})
		}

		// Generate the public URL for the image
		publicUrl := fmt.Sprintf("/images/%s", file.Filename)

		usersStr := c.FormValue("users")
		users, err := strconv.Atoi(usersStr) // Convert to integer
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid number of users",
			})
		}

		// Create a new product
		product := &models.Product{
			Name:        c.FormValue("name"),
			Description: c.FormValue("description"),
			Email:       c.FormValue("email"),
			Image:       publicUrl, // Store the public URL in the database
			Users:       users,
			Url:         c.FormValue("url"),
		}
		fmt.Println(product)
		database.DB.DB.Create(product)

		products := models.Product{}

		// Find the most recent product by ordering by ID descending and getting the first
		result := database.DB.DB.Order("id desc").First(&products)
		if result.Error != nil {
			// If there's an error (e.g., no products found), return a 404 status
			return c.Status(fiber.StatusNotFound).SendString("No products found")
		}

		// Return the last product
		return c.Render("projects", fiber.Map{
			"data": []models.Product{products}, // Wrap in a slice for consistent rendering
		})
	})

	app.Get("/", func(c *fiber.Ctx) error {
		product := models.Product{}

		// Find the most recent product by ordering by ID descending and getting the first
		result := database.DB.DB.Order("id desc").First(&product)
		if result.Error != nil {
			// If there's an error (e.g., no products found), return a 404 status
			return c.Status(fiber.StatusNotFound).SendString("No products found")
		}

		// Return the last product
		return c.Render("index", fiber.Map{
			"data": []models.Product{product}, // Wrap in a slice for consistent rendering
		})
	})

	app.Delete("/delete", func(c *fiber.Ctx) error {
		// Perform a DELETE operation on the products table
		result := database.DB.DB.Where("1 = 1").Delete(&models.Product{}) // Deletes all rows in the table

		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to delete products",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "All products deleted successfully",
		})
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Heelo Ji")
	})

	app.Static("/images", "./public/images")

	log.Fatal(app.Listen(":3000"))
}
