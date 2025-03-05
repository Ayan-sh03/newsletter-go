package routes

import (
	"net/http"
	controllers "newsletter/internal/api/controllers/newsletter"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// main routes
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))

	}).Methods("GET")

	// api routes

	r.HandleFunc("/api/v1/newsletter", controllers.CreateNewsletter).Methods("POST")

	return r
}
