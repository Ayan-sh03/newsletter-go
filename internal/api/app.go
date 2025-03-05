package api

import (
	"fmt"
	"log"
	"net/http"
	"newsletter/internal/routes"
	"newsletter/internal/storage/db"

	"github.com/gorilla/handlers"
)

func Run() {

	router := routes.SetupRoutes()
	db.InitDb()
	port := ":8000"

	// CORS middleware
	headersOk := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// Wrap your router with the CORS middleware
	corsRouter := handlers.CORS(headersOk, originsOk, methodsOk)(router)
	fmt.Println("Server started on port :" + port)
	log.Fatal(http.ListenAndServe(port, corsRouter))

}
