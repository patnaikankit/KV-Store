package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/patnaikankit/KV-Store.git/pkg/controllers"
	"github.com/patnaikankit/KV-Store.git/pkg/store"
)

const PORT = 4000

func main() {
	store.InitKvStore("../kv-data.json")

	mux := http.NewServeMux()

	mux.HandleFunc("/get", controllers.GetKey)
	mux.HandleFunc("/set", controllers.SetKey)
	mux.HandleFunc("/update", controllers.UpdateKey)
	mux.HandleFunc("/delete", controllers.DeleteKey)

	server := fmt.Sprintf(":%d", PORT)
	fmt.Printf("Starting server on port: %d\n", PORT)

	err := http.ListenAndServe(server, mux)
	if err != nil {
		log.Fatalf("Error while starting server: %v", err)
	}
}
