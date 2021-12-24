package main

import (
	"main/app/controllers"
	"main/app/models"
)

func main() {

	models.ConnectionDB()
	models.TestDb()
	controllers.StartWebServer()
}
