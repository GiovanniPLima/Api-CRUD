package main

import (
	controller "crud/Controllers"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	//Cria a Rota
	router := mux.NewRouter()
	//CREATE - POST
	router.HandleFunc("/usuarios", controller.CriarUsuario).Methods(http.MethodPost)
	//READ - GET
	router.HandleFunc("/usuarios", controller.BuscarUsuarios).Methods(http.MethodGet)
	//READ - GET BY ID
	router.HandleFunc("/usuarios/{id}", controller.BuscarUsuario).Methods(http.MethodGet)

	// Cria o servidor Par fazer o C-R-U-D
	fmt.Println("Servidor Rodando")
	log.Fatal(http.ListenAndServe(":5000", router))
}
