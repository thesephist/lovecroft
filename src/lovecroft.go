package main

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func index(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok")
}

func subscribe(w http.ResponseWriter, r *http.Request) {

}

func unsubscribe(w http.ResponseWriter, r *http.Request) {

}

func directory(w http.ResponseWriter, r *http.Request) {

}

func list(w http.ResponseWriter, r *http.Request) {

}

func start() {
	r := mux.NewRouter()

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:7171",
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}

	r.HandleFunc("/", index)

	r.Methods("POST").Path("/subscribe").HandlerFunc(subscribe)
	r.Methods("GET").Path("/unsubscribe").HandlerFunc(unsubscribe)

	r.Methods("GET").Path("/directory").HandlerFunc(directory)
	r.Methods("GET").Path("/list/{listName}").HandlerFunc(list)

	log.Printf("Lovecroft listening on %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func main() {
	log.Println("Starting Lovecroft server")

	start()
}
