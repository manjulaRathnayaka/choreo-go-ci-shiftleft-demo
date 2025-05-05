package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() (*gin.Engine, *Store) {
	store := &Store{
		todos:  []Todo{},
		nextID: 1,
	}
	router := gin.Default()
	router.POST("/todos", store.createTodo)
	router.GET("/todos", store.listTodos)
	return router, store
}

func TestCreateTodo_Success(t *testing.T) {
	router, _ := setupRouter()
	payload := map[string]string{"title": "Test Todo"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp Todo
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Test Todo", resp.Title)
	assert.False(t, resp.Done)
	assert.Equal(t, 1, resp.ID)
}

func TestCreateTodo_ValidationFail(t *testing.T) {
	router, _ := setupRouter()
	payload := map[string]string{"title": ""}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListTodos(t *testing.T) {
	router, store := setupRouter()

	// Seed data
	store.Lock()
	store.todos = append(store.todos, Todo{ID: 1, Title: "Sample", Done: false})
	store.Unlock()

	req, _ := http.NewRequest("GET", "/todos", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []Todo
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "Sample", resp[0].Title)
}
