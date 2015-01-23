package main

import (
	"encoding/json"
	"fmt"
	"github.com/albrow/go-data-parser"
	"github.com/albrow/negroni-json-recovery"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/martini-contrib/cors"
	"github.com/unrolled/render"
	"net/http"
	"sort"
	"strconv"
	"sync"
)

type Todo struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	IsCompleted bool   `json:"isCompleted"`
}

type todosIndex map[int]*Todo
type todosList []*Todo

var (
	// todos stores all the todos as a map of id to *Todo
	todos = todosIndex{}
	// todosMutex protects access to the todos map
	todosMutex = sync.Mutex{}
	// todosCounter is incremented every time a new todo is created
	// it is used to set todo ids.
	todosCounter = 0
)

const (
	statusUnprocessableEntity = 422
)

func main() {
	createInitialTodos()

	// Routes
	router := mux.NewRouter()
	router.HandleFunc("/todos", todosController.Index).Methods("GET")
	router.HandleFunc("/todos", todosController.Create).Methods("POST")
	router.HandleFunc("/todos/{id}", todosController.Update).Methods("PUT")
	router.HandleFunc("/todos/{id}", todosController.Delete).Methods("DELETE")

	// Other middleware
	n := negroni.New(negroni.NewLogger())
	n.UseHandler(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	n.Use(recovery.JSONRecovery(true))
	recovery.StackDepth = 3

	// Router must always come last
	n.UseHandler(router)

	// Start the server
	n.Run(":3000")
}

func createInitialTodos() {
	createTodo("Write a frontend framework in Go")
	createTodo("???")
	createTodo("Profit!")
}

func createTodo(title string) *Todo {
	todosMutex.Lock()
	defer todosMutex.Unlock()
	id := todosCounter
	todosCounter++
	todo := &Todo{
		Id:    id,
		Title: title,
	}
	todos[id] = todo
	return todo
}

// Todos Controller and its methods
type todosControllerType struct{}

var todosController = todosControllerType{}

func (todosControllerType) Index(w http.ResponseWriter, req *http.Request) {
	r := render.New()
	r.JSON(w, http.StatusOK, todos)
}

func (todosControllerType) Create(w http.ResponseWriter, req *http.Request) {
	r := render.New()

	// Parse data and do validations
	todoData, err := data.Parse(req)
	if err != nil {
		panic(err)
	}
	val := todoData.Validator()
	val.Require("title")
	if val.HasErrors() {
		r.JSON(w, statusUnprocessableEntity, val.ErrorMap())
		return
	}

	// Create the todo and render response
	todo := createTodo(todoData.Get("title"))
	r.JSON(w, http.StatusOK, todo)
}

func (todosControllerType) Update(w http.ResponseWriter, req *http.Request) {
	r := render.New()

	// Get the existing todo from the map or render an error
	// if it wasn't found
	urlParams := mux.Vars(req)
	idString := urlParams["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		panic(err)
	}
	todo, found := todos[id]
	if !found {
		msg := fmt.Sprintf("Could not find todo with id = %d", id)
		r.JSON(w, http.StatusNotFound, map[string]string{
			"error": msg,
		})
		return
	}

	// Update the todo with the data in the request
	todoData, err := data.Parse(req)
	if err != nil {
		panic(err)
	}
	todosMutex.Lock()
	if todoData.KeyExists("title") {
		todo.Title = todoData.Get("title")
	}
	if todoData.KeyExists("isCompleted") {
		todo.IsCompleted = todoData.GetBool("isCompleted")
	}
	todosMutex.Unlock()

	// Render response
	r.JSON(w, http.StatusOK, todo)
}

func (todosControllerType) Delete(w http.ResponseWriter, req *http.Request) {
	r := render.New()

	// Get the id from the url parameters
	urlParams := mux.Vars(req)
	idString := urlParams["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		panic(err)
	}

	// Delete the todo and render a response
	todosMutex.Lock()
	delete(todos, id)
	todosMutex.Unlock()
	r.JSON(w, http.StatusOK, struct{}{})
}

// Make todosIndex satisfy the json.Marshaller interface
// It will return a json array of todos sorted by id
func (t todosIndex) MarshalJSON() ([]byte, error) {
	todosList := todosList{}
	for _, todo := range t {
		todosList = append(todosList, todo)
	}
	sort.Sort(todosList)
	return json.Marshal(todosList)
}

// Make todoList satisfy sort.Interface
func (tl todosList) Len() int {
	return len(tl)
}

func (tl todosList) Less(i, j int) bool {
	return tl[i].Id < tl[j].Id
}

func (tl todosList) Swap(i, j int) {
	tl[i], tl[j] = tl[j], tl[i]
}
