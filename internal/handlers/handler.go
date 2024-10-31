package handlers

import (
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler chamado!")

	if _, err := w.Write([]byte("Hello, world!")); err != nil {
		log.Printf("Erro ao processar a requisição: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
