package main

import (
	"database/sql"
	"log"
	"radiusgo/controllers"
	"radiusgo/middlewares"
	"radiusgo/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"
)

var jwtSecret = []byte("chavemtsecreta")

var db *sql.DB

func main() {
	conn := "postgres://postgres:12345@localhost:5432/radius?sslmode=disable"
	var err error
	db, err = sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}

	services.AutoClean(db, 24*time.Hour)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	app.Post("/login", controllers.HandleLogin(db))
	app.Post("/refresh", controllers.HandleRefreshToken(db))
	app.Post("/register", middlewares.AuthMiddleware, middlewares.RoleMiddleware("admin"), controllers.HandleRegister(db))
	app.Get("/users", middlewares.AuthMiddleware, middlewares.RoleMiddleware("user", "admin"), controllers.HandleGetUsers(db))
	app.Post("/user/password", middlewares.AuthMiddleware, middlewares.RoleMiddleware("user", "admin"), controllers.HandlePasswordChange(db))

	app.Post("/bilhetes", controllers.HandlePostBilhete(db))
	app.Get("/bilhetes", middlewares.AuthMiddleware, middlewares.RoleMiddleware("user", "admin"), controllers.HandleGetBilhetes(db))

	app.Post("/rotinas/limpezadias", middlewares.AuthMiddleware, middlewares.RoleMiddleware("admin"), controllers.HandleUpdateCleanupDays(db))
	app.Get("/rotinas/limpezadias", middlewares.AuthMiddleware, middlewares.RoleMiddleware("admin"), controllers.HandleGetCleanupDays(db))

	app.Get("/gateways", middlewares.AuthMiddleware, middlewares.RoleMiddleware("admin"), controllers.HandleGetGateways(db))
	app.Post("/gateways", middlewares.AuthMiddleware, middlewares.RoleMiddleware("admin"), controllers.HandlePostGateways(db))
	app.Delete("/gateways", middlewares.AuthMiddleware, middlewares.RoleMiddleware("admin"), controllers.HandleDeleteGateways(db))
	log.Fatal(app.Listen(":5000"))
}
