package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"main.go/database"
	"main.go/models"
)

func main() {

	database.ConnectDb()

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})

	fmt.Println("Jai Shree Ram !!")

	imageDir := "./public/images/"
	if _, err := os.Stat(imageDir); os.IsNotExist(err) {
		os.MkdirAll(imageDir, 0755) // Create the directory if it doesn't exist
	}

	app.Post("/new", func(c *fiber.Ctx) error {
		fmt.Println(c.FormFile("image"))
		file, err := c.FormFile("image")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "No file uploaded",
			})
		}
		fmt.Println(file)
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
			Name:    c.FormValue("name"),
			Description: c.FormValue("description"),
			Email:   c.FormValue("email"),
			Image: publicUrl, // Store the public URL in the database
			Users: users,
		}
		database.DB.DB.Create(product)

		products := []models.Product{}
		database.DB.DB.Find(&products)

		// Render the updated list of products
		return c.Render("projects", fiber.Map{
			"data": products,
		})
	})

	app.Get("/", func(c *fiber.Ctx) error {
		products := []models.Product{}

		database.DB.DB.Find(&products)
		fmt.Println(products)
		return c.Render("index",fiber.Map{
			"data": products,
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

	app.Static("/images", "./public/images")

	log.Fatal(app.Listen(":3000"))
}
