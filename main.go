package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
	"template/database"
	"template/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func init() {
	database.ConnectDb() // for connect Db define in function.go
}

func main() {

	// Define custom functions
	funcMap := template.FuncMap{
		"sub": func(a, b int) int {
			return a - b
		},
		"add": func(a, b int) int {
			return a + b
		},
		"mul": func(a, b int) int {
			return a * b
		},
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
	}

	// Load templates with custom functions
	engine := html.New("./views", ".html")
	engine.AddFuncMap(funcMap)

	// Register the function `tolower`
	engine.AddFunc("tolower", strings.ToLower)

	// Register the function `toupper`
	engine.AddFunc("toupper", strings.ToUpper)

	app := fiber.New(fiber.Config{
		// Sets the view engine for rendering templates (e.g., HTML).
		Views: engine,
		// Specifies the default layout template to be used for views.
		ViewsLayout: "layouts/main",
		// Passes local variables defined in Fiber handlers to the view templates automatically.
		PassLocalsToViews: true,
	})
	// Apply CORS middleware globally
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",                            // Allow all origins or specify certain domains
		AllowMethods:     "GET, POST, PATCH, DELETE",     // Allow only GET requests
		AllowHeaders:     "Origin, Content-Type, Accept", // Allow necessary headers
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	}))
	app.Static("/views", "./views") // allow folder for get public files like js/css etc

	// Serve static files from the "static" directory with MIME type
	app.Static("/assets", "./assets", fiber.Static{
		Compress:      true,
		CacheDuration: 10 * 60 * 1000, // Cache static files
		MaxAge:        3600,
		ByteRange:     true,
		Browse:        true,
		Download:      false,
	})
	// Middleware to set correct Content-Type for CSS files
	app.Use("/assets/css", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/css")
		return c.Next()
	})

	//app.Static("/assets", "./assets") // allow folder for get public files like js/css etc
	app.Use(logger.New())

	// check .env file is exist or not
	if err := godotenv.Load(".env"); err != nil {
		fmt.Printf("ENV not Found")
		return
	}

	// Middleware to set BaseURL globally
	app.Use(func(c *fiber.Ctx) error {
		// Set the BaseURL using ctx.Locals
		c.Locals("CommonURL", os.Getenv("CommonURL"))     // Set CommonURL File Path
		c.Locals("LogoLight", os.Getenv("LogoLight"))     // Set Logo Path
		c.Locals("LogoDark", os.Getenv("LogoDark"))       // Set Logo Path
		c.Locals("FaviconIcon", os.Getenv("FaviconIcon")) // Set FaviconIcon Path
		c.Locals("HostName", os.Getenv("HostName"))       // Set FaviconIcon Path

		GetURL := c.BaseURL() // Get Base Url
		if GetURL == "http://localhost:"+os.Getenv("PORT") {
			c.Locals("CssURLS", os.Getenv("GitUrl")) // Set Logo Path
		} else {
			c.Locals("CssURLS", os.Getenv("FileURL")) // Set Logo Path
		}
		return c.Next()
	})

	routes.InitRoutes(app) // for Set page routes path define in routes.go

	// Handle 404 errors with a custom HTML page
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).Render("404", fiber.Map{})
	})

	// Starts the Fiber application and listens on the port specified in the "PORT" environment variable.
	// Logs a fatal error if the application fails to start.
	log.Fatal(app.Listen(":" + os.Getenv("PORT") + ""))

}
