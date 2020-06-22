package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type DirectoryStore struct {
	root      string
	directory Directory
}

func (ds *DirectoryStore) InstantiateDirectory() {
	err := os.MkdirAll(ds.root, 0755)
	if err != nil {
		panic("DirectoryStore.EnsureExists: " + err.Error())
	}

	entries, err := ioutil.ReadDir(ds.root)
	if err != nil {
		panic("DirectoryStore.EnsureExists: " + err.Error())
	}

	ds.directory.Lists = make([]List, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".csv") {
			continue
		}

		file, err := os.Open(filepath.Join(ds.root, name))
		if err != nil {
			panic("Error reading database: " + err.Error())
		}
		defer file.Close()

		scribers := []Subscriber{}
		reader := csv.NewReader(file)
		for {
			line, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				panic("Error reading database CSV: " + err.Error())
			}

			startDate, err := time.Parse(time.RFC3339, line[3])
			if err != nil {
				panic(fmt.Sprintf("Error reading database %s: %s", name, err.Error()))
			}
			endDate, err := time.Parse(time.RFC3339, line[4])
			if err != nil {
				panic(fmt.Sprintf("Error reading database %s: %s", name, err.Error()))
			}

			scribers = append(scribers, Subscriber{
				Email:      line[0],
				GivenName:  line[1],
				FamilyName: line[2],
				StartDate:  startDate,
				EndDate:    endDate,
				UnsubToken: line[5],
			})
		}

		ds.directory.Lists = append(ds.directory.Lists, List{
			Name:        name[0 : len(name)-4],
			Subscribers: scribers,
		})
	}

	log.Printf("Directory initialized with %d lists:\n", len(ds.directory.Lists))
	for _, list := range ds.directory.Lists {
		log.Printf("\t%s - %d actives\n", list.Name, len(list.ActiveSubscribers()))
	}
}

func (ds *DirectoryStore) Commit() error {
	for _, list := range ds.directory.Lists {
		file, err := os.OpenFile(
			filepath.Join(ds.root, fmt.Sprintf("%s.csv", list.Name)),
			os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
			0644,
		)
		if err != nil {
			return err
		}
		defer file.Close()

		io.WriteString(file, list.RenderToCSV())
	}

	return nil
}
