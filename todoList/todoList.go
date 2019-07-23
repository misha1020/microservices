package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type todo struct {
	Key   int
	Value string
}

var (
	id    int
	todos map[int]string
)

func main() {
	todos = make(map[int]string)
	r := mux.NewRouter()
	r.HandleFunc("/", CreateHandler).Methods("POST")
	r.HandleFunc("/{id:[0-9]+}", UpdateHandler).Methods("PUT")
	r.HandleFunc("/{id:[0-9]+}", DeleteHandler).Methods("DELETE")
	r.HandleFunc("/", GetAllHandlers).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}", GetHandler).Methods("GET")
	r.HandleFunc("/ctx", ContextHandler).Methods("GET")
	log.Print("Server has started")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Print(err)
	}
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	//var todo todo
	//err := getItem(r, &todo)
	todo, err := getItem(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("Status bad request")
		return
	}

	todos[id] = todo.Value
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("todo[" + strconv.Itoa(id) + "] has been added"))
	log.Printf("%s has been added with id %d", todo.Value, id)
	id++
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, _ := strconv.Atoi(idStr)
	_, ok := todos[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("todo[%d] does not exist", id)
		return
	}

	todo, err := getItem(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("Status bad request")
		return
	}
	todos[id] = todo.Value
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("todo[" + strconv.Itoa(id) + "] has been updated"))
	log.Printf("todo[%d] updated value to %s", id, todo.Value)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, _ := strconv.Atoi(idStr)
	_, ok := todos[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("todo[%d] does not exist", id)
		return
	}
	delete(todos, id)
	w.WriteHeader(http.StatusNoContent)
	_, _ = w.Write([]byte("todo[" + strconv.Itoa(id) + "] has been deleted"))
	log.Printf("todo[%d] has been deleted", id)
}

func GetAllHandlers(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	log.Print("All values have been printed")
	for k, v := range todos {
		_, _ = w.Write([]byte(strconv.Itoa(k)))
		_, _ = w.Write([]byte(": " + v + "\n"))
	}
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, _ := strconv.Atoi(idStr)
	_, ok := todos[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("todo[%d] does not exist", id)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("todo[" + strconv.Itoa(id) + "] = " + todos[id]))
	log.Printf("todo[%d] has been printed", id)
}

func ContextHandler(w http.ResponseWriter, r *http.Request) {
	context := r.Context()

	select {
	case <-time.After(5 * time.Second):
		log.Print("all ok")
		w.WriteHeader(http.StatusOK)
	case <-context.Done():
		log.Print("current operation was cancelled")
	}
}

func getItem(r *http.Request) (*todo, error) {
	var todo *todo
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &todo)
	if err != nil {
		return nil, err
	}

	return todo, nil
}