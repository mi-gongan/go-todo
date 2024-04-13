package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTodos(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler())
	defer ts.Close()

	resp, err := http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test todo"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	var  todo Todo
	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal("Test todo", todo.Name)
	id1 := todo.Id

	resp, err = http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test todo 2"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal("Test todo 2", todo.Name)
	id2 := todo.Id
	
	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos := []*Todo{}
	json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(2, len(todos))
	for _, t := range todos {
		if t.Id == id1 {
			assert.Equal("Test todo", t.Name)
		} else if t.Id == id2 {
			assert.Equal("Test todo 2", t.Name)
		} else {
			assert.Error(fmt.Errorf("Unexpected todo: %v", t))
		}
	}

	resp,err = http.Get(ts.URL + "/complete-todo/" + strconv.Itoa(id1) +"?complete=true")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos = []*Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(2, len(todos))
	for _, t := range todos {
		if t.Id == id1 {
			assert.True(t.Completed)
		} 
	}

}