package api

import (
	"ecommerce-golang/config"
	"ecommerce-golang/internal/api/rest"
	"ecommerce-golang/internal/api/rest/handlers"
	"ecommerce-golang/internal/domain"
	"ecommerce-golang/internal/helper"
	"log"

	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func StartServer(config config.AppConfig){
	app := fiber.New();

	db, err := gorm.Open(postgres.Open(config.Dsn), &gorm.Config{});

	if err != nil{
		log.Fatalf("database connection error %v\n", err)
	}

	log.Printf("database connected")
	
	err = db.AutoMigrate(&domain.User{}, &domain.BankAccount{}, &domain.Category{},&domain.Product{})
	if err != nil {
		log.Fatalf("error on running migration %v", err.Error()) 
	}
	log.Println("migration successful") 
	auth := helper.SetupAuth(config.AppSecret)

	rh := &rest.RestHandler{
		App: app,
		DB: db,
		Auth: auth,
		Config: config,
	}

	setupRoutes(rh)

	app.Listen(config.ServerPort)

}	

func setupRoutes(rh *rest.RestHandler) {
	// user handler
	handlers.SetupUserRoutes(rh)
	// transactions
	// catalog
	handlers.SetupCatalogRoutes(rh)
}