package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Task struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

var tasks []Task

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/tasks", handleTasksHandler).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	r.HandleFunc("/tasks/{id}", handleTaskById).Methods(http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodOptions)
	r.Use(mux.CORSMethodMiddleware(r))

	http.ListenAndServe(":8080", r)
}

// To handle GET /tasks and POST /tasks request
func handleTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	w.Header().Set("Content-Type", "application/json")

	var data []byte
	var err error

	// Creates new tasks and store it in global state variable
	if r.Method == http.MethodPost {
		// Parsing task name from request body
		var new_task Task
		err = json.NewDecoder(r.Body).Decode(&new_task)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Sorry! An error occurred"))
			return
		}

		// Generates a unique key for tasks id
		id := uuid.New()
		new_task.Id = id.String()

		// Appends data to global variable
		tasks = append(tasks, new_task)

		data, err = json.Marshal(new_task)
	} else if r.Method == http.MethodGet { // Returns all tasks from the global state variable
		data, err = json.Marshal(tasks)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found!"))
		return
	}

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Sorry! An error occurred"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleTaskById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var data []byte
	var err error
	var task Task

	vars := mux.Vars(r)
	task_id := vars["id"]

	if r.Method == http.MethodGet { // To handle get task by id operation
		task, _ = getTaskById(task_id)
		data, err = json.Marshal(task)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Sorry! An error occurred"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	} else if r.Method == http.MethodPut { // To handle update task operation
		_, index := getTaskById(task_id)

		var new_task Task
		err = json.NewDecoder(r.Body).Decode(&new_task)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Sorry! An error occurred"))
			return
		}

		tasks[index].Name = new_task.Name

		data, err = json.Marshal(tasks[index])
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Sorry! An error occurred"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	} else if r.Method == http.MethodDelete { // To handle delete task operation
		tasks = slices.DeleteFunc(tasks, func(item Task) bool {
			return item.Id == task_id
		})
		w.WriteHeader(http.StatusNoContent)
	}
}

func getTaskById(task_id string) (Task, int) {
	for index, item := range tasks {
		if item.Id == task_id {
			return item, index
		}
	}

	return Task{}, 0
}
