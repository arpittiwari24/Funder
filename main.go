package main

import (
	"fmt"
	"log"

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

	app.Post("/new", func(c *fiber.Ctx) error {
		product := new(models.Product)

		if err := c.BodyParser(product) ; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message" : err.Error(),
			})
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

	log.Fatal(app.Listen(":3000"))
}
