package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
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

func makeSubscribe(ds *DirectoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		list, err := ds.directory.FindList(vars["listName"])
		if err != nil {
			sendError(w, err)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			sendError(w, err)
			return
		}

		scriber := Subscriber{}
		err = json.Unmarshal(body, &scriber)
		if err != nil {
			sendError(w, err)
			return
		}

		list.Subscribe(scriber)
		err = ds.Commit()
		if err != nil {
			sendError(w, err)
		}
	}
}

func makeUnsubscribe(ds *DirectoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		list, err := ds.directory.FindList(vars["listName"])
		if err != nil {
			sendError(w, err)
			return
		}

		err = list.Unsubscribe(vars["token"])
		if err != nil {
			sendError(w, err)
			return
		}
		err = ds.Commit()
		if err != nil {
			sendError(w, err)
		}
	}
}

func makeDirectory(ds *DirectoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := useTemplate("directory")
		if err != nil {
			sendError(w, err)
			return
		}

		err = tmpl.Execute(w, ds.directory)
		if err != nil {
			sendError(w, err)
			return
		}
	}
}

func makeList(ds *DirectoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := useTemplate("list")
		if err != nil {
			sendError(w, err)
			return
		}

		vars := mux.Vars(r)
		list, err := ds.directory.FindList(vars["listName"])
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

func makeListCSV(ds *DirectoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		list, err := ds.directory.FindList(vars["listName"])
		if err != nil {
			sendError(w, err)
			return
		}

		io.WriteString(w, list.RenderToCSV())
	}
}

func start() {
	store := &DirectoryStore{
		root: "./db/",
	}
	store.InstantiateDirectory()

	r := mux.NewRouter()

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:7171",
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method, r.RequestURI)

			next.ServeHTTP(w, r)
		})
	})

	r.HandleFunc("/", index)

	r.Methods("POST").Path("/subscribe/{listName}").HandlerFunc(makeSubscribe(store))
	r.Methods("GET").Path("/unsubscribe/{listName}/{token}").HandlerFunc(makeUnsubscribe(store))

	r.Methods("GET").Path("/directory").HandlerFunc(makeDirectory(store))
	r.Methods("GET").Path("/list/{listName}").HandlerFunc(makeList(store))
	r.Methods("GET").Path("/list-csv/{listName}").HandlerFunc(makeListCSV(store))

	log.Printf("Lovecroft listening on %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func main() {
	log.Println("Starting Lovecroft server")

	start()
}
