package main

import (
	"fmt"

	"github.com/Romasmi/go-rest-api-template/internal/application"
)

// @title Go REST API Template
// @version 1.0
// @description A RESTful API template using Go
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	app := &application.App{}
	err := app.InitApp("../../")
	if err != nil {
		fmt.Printf("error while app initialization: %v", err)
		return
	}
	defer app.OnStop()

	app.Run()
}
