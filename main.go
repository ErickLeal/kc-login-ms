package main

import (
	"log"
	"net/http"

	"login-ms/handlers"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Login Microservice is running"))
	})

	http.HandleFunc("/create-client", handlers.CreateClientHandler)
	http.HandleFunc("/generate-pkce", handlers.GeneratePKCEHandler)
	http.HandleFunc("/redirect", handlers.RedirectHandler)

	log.Println("Starting server on :8005")
	if err := http.ListenAndServe(":8005", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
