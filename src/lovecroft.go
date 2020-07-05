package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func sendError(w http.ResponseWriter, err error) {
	log.Println(err.Error())
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, err.Error())
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := useTemplate("index")
	if err != nil {
		sendError(w, err)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		sendError(w, err)
		return
	}
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
			return
		}

		// allow subscriptions to come in from any site
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
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
			return
		}

		tmpl, err := useTemplate("unsubscribe")
		if err != nil {
			sendError(w, err)
			return
		}

		scriber, err := list.SubscriberFromToken(vars["token"])
		if err != nil {
			sendError(w, err)
			return
		}

		err = tmpl.Execute(w, scriber)
		if err != nil {
			sendError(w, err)
			return
		}
	}
}

func makeCreateList(ds *DirectoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		listName := vars["listName"]
		for _, list := range ds.directory.Lists {
			if list.Name == listName {
				w.WriteHeader(http.StatusConflict)
				return
			}
		}

		ds.directory.Lists = append(ds.directory.Lists, List{
			Name:        listName,
			Subscribers: []Subscriber{},
		})

		err := ds.Commit()
		if err != nil {
			sendError(w, err)
			return
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

func makeAuthor(ds *DirectoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := useTemplate("author")
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

func makeSend(ds *DirectoryStore) http.HandlerFunc {
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

		issue := Issue{}
		err = json.Unmarshal(body, &issue)
		if err != nil {
			sendError(w, err)
			return
		}

		go list.Send(issue)
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

func makeListActiveCSV(ds *DirectoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		list, err := ds.directory.FindList(vars["listName"])
		if err != nil {
			sendError(w, err)
			return
		}

		io.WriteString(w, list.RenderActivesToCSV())
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
			if !strings.HasPrefix(r.RequestURI, "/static") {
				log.Println(r.Method, r.RequestURI)
			}

			next.ServeHTTP(w, r)
		})
	})

	r.HandleFunc("/", index)

	// public paths
	r.Methods("POST").Path("/subscribe/{listName}").HandlerFunc(makeSubscribe(store))
	r.Methods("GET").Path("/unsubscribe/{listName}/{token}").HandlerFunc(makeUnsubscribe(store))

	// admin paths
	r.Methods("POST").Path("/admin/create-list/{listName}").HandlerFunc(makeCreateList(store))
	r.Methods("GET").Path("/admin/directory").HandlerFunc(makeDirectory(store))
	r.Methods("GET").Path("/admin/author").HandlerFunc(makeAuthor(store))
	r.Methods("GET").Path("/admin/list/{listName}").HandlerFunc(makeList(store))
	r.Methods("GET").Path("/admin/list-csv/{listName}.csv").HandlerFunc(makeListCSV(store))
	r.Methods("GET").Path("/admin/list-active-csv/{listName}.csv").HandlerFunc(makeListActiveCSV(store))
	r.Methods("POST").Path("/admin/send/{listName}").HandlerFunc(makeSend(store))

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Printf("Lovecroft listening on %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func main() {
	log.Println("Starting Lovecroft server")

	start()
}
