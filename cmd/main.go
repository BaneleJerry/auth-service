package main

import (
	"log"
	"net/http"

	"github.com/BaneleJerry/auth-service/app"
	"github.com/BaneleJerry/auth-service/database"
	"github.com/BaneleJerry/auth-service/routes"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := database.ConnectDB()

	application := app.InitApp(db)

	router := routes.SetupRouter(&application.Controllers)

	log.Fatal(http.ListenAndServe(":8080", router))

}
