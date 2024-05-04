package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"main.go/database"
)

func main() {

	database.ConnectDb()

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	fmt.Println("Jai Shree Ram !!")

	app.Post("/new", func(c *fiber.Ctx) error {
		return nil
	})

	log.Fatal(app.Listen(":3000"))
}
