package main

import (
	"fmt"
	"lapas/pkg/controllers"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	origins := handlers.AllowedOrigins([]string{"*"})

	router.HandleFunc("/user/{idUser}", controllers.GetUser).Methods("GET")
	router.HandleFunc("/user", controllers.GetUsers).Methods("GET")
	router.HandleFunc("/user", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/user/{idUser}", controllers.UpdateUser).Methods("PUT")
	router.HandleFunc("/password", controllers.ResetPassword).Methods("POST")
	router.HandleFunc("/password/{idUser}", controllers.ChangePassword).Methods("PUT")

	fmt.Println("Server running at: 5000")
	log.Fatal(http.ListenAndServe(":5000", handlers.CORS(origins)(router)))

}
