package main

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func sendError(w http.ResponseWriter, err error) {
	log.Println(err.Error())
	io.WriteString(w, err.Error())
}

func index(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok")
}

func makeSubscribe(directory *Directory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func makeUnsubscribe(directory *Directory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func makeDirectory(directory *Directory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := useTemplate("directory")
		if err != nil {
			sendError(w, err)
			return
		}

		err = tmpl.Execute(w, directory)
		if err != nil {
			sendError(w, err)
			return
		}
	}
}

func makeList(directory *Directory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := useTemplate("list")
		if err != nil {
			sendError(w, err)
			return
		}

		vars := mux.Vars(r)
		list, err := directory.FindList(vars["listName"])
		if err != nil {
			sendError(w, err)
			return
		}

		err = tmpl.Execute(w, list)
		if err != nil {
			sendError(w, err)
			return
		}
	}
}

func makeListCSV(directory *Directory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		list, err := directory.FindList(vars["listName"])
		if err != nil {
			sendError(w, err)
			return
		}

		io.WriteString(w, list.RenderToCSV())
	}
}

func start() {
	store := DirectoryStore{
		root: "./db/",
	}
	dir := store.InstantiateDirectory()

	r := mux.NewRouter()

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:7171",
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}

	r.HandleFunc("/", index)

	r.Methods("POST").Path("/subscribe").HandlerFunc(makeSubscribe(dir))
	r.Methods("GET").Path("/unsubscribe").HandlerFunc(makeUnsubscribe(dir))

	r.Methods("GET").Path("/directory").HandlerFunc(makeDirectory(dir))
	r.Methods("GET").Path("/list/{listName}").HandlerFunc(makeList(dir))
	r.Methods("GET").Path("/list-csv/{listName}").HandlerFunc(makeListCSV(dir))

	log.Printf("Lovecroft listening on %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func main() {
	log.Println("Starting Lovecroft server")

	start()
}
