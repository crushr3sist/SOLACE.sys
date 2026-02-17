package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/template/html/v3"
)

//go:embed frontend/*
var frontend_dir embed.FS

//go:embed frontend/static/*
var static_dir embed.FS

func main() {

	// create a views engine
	engine := html.NewFileSystem(http.FS(frontend_dir), ".html")

	// start a new fiber instance
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// embed the static files
	static_files := fs.FS(static_dir)

	// use those static files
	app.Use("/", static.New("", static.Config{FS: static_files}))

	// serve the index.html page
	app.Get("/", func(c fiber.Ctx) error {
		return c.Render("frontend/index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	fmt.Println("listening on http://localhost:3000/")

	// start the server
	log.Fatal(app.Listen(":3000"))

}
