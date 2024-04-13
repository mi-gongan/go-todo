package app

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todos", http.StatusTemporaryRedirect)
}

var rd *render.Render

type Todo struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Completed bool `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type Success struct {
	Success bool `json:"success"`
}

var todoMap map[int]*Todo

func getTodoListHandler(w http.ResponseWriter, r *http.Request) {
	list := []*Todo{}
	for _, v := range todoMap {
		list = append(list, v)
	}
	rd.JSON(w, http.StatusOK, list)
}

func addTodoHandler(w http.ResponseWriter, r *http.Request) {
	t := &Todo{}
	t.Id = len(todoMap) + 1
	t.Name = r.FormValue("name")
	t.Completed = false
	t.CreatedAt = time.Now()

	todoMap[t.Id] = t

	rd.JSON(w, http.StatusCreated, t)
}

func removeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	if _, ok := todoMap[id]; ok {
		delete(todoMap, id)
		rd.JSON(w, http.StatusOK, Success{true})
		return
	} else {
		rd.JSON(w, http.StatusNotFound, Success{false})
		return
	}
}

func completeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	complete := r.FormValue("complete") == "true"
	if todo, ok := todoMap[id]; ok {
		todo.Completed = complete
		rd.JSON(w, http.StatusOK, Success{true})
		return
	} else {
		rd.JSON(w, http.StatusNotFound, Success{false})
		return
	}

}

func addTestTodos() {
	todoMap[1] = &Todo{1, "Buy a milk", false, time.Now()}
	todoMap[2] = &Todo{2, "Extercise", true, time.Now()}
	todoMap[3] = &Todo{3, "Learn Go", false, time.Now()}


}

func MakeHandler() http.Handler {
	todoMap = make(map[int]*Todo)
	addTestTodos()
	rd = render.New()
	r := mux.NewRouter()

	r.HandleFunc("/todos", getTodoListHandler).Methods("GET")
	r.HandleFunc("/todos", addTodoHandler).Methods("POST")
	r.HandleFunc("/todos/{id:[0-9]+}", removeTodoHandler).Methods("DELETE")
	r.HandleFunc("/complete-todo/{id:[0-9]+}", completeTodoHandler).Methods("GET")

	r.HandleFunc("/", indexHandler)
	
	return r
}


