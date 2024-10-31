package handlers

import (
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Handler chamado!")); err != nil {
		log.Printf("Erro ao processar a requisição: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
