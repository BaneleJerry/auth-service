package main

import (
	"log"
	"net/http"

	"github.com/BaneleJerry/auth-service/app"
	"github.com/BaneleJerry/auth-service/database"
	"github.com/BaneleJerry/auth-service/routes"
)

func main() {

	db := database.ConnectDB()

	application := app.InitApp(db)

	router := routes.SetupRouter(&application.Controllers)

	log.Fatal(http.ListenAndServe(":8080", router))

}
