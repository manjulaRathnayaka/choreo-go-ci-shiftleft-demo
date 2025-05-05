package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// Todo represents a task with a title and completion status.
type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// Store holds the in-memory list of todos.
type Store struct {
	sync.Mutex
	todos  []Todo
	nextID int
}

func main() {
	router := gin.Default()
	store := &Store{
		todos:  []Todo{},
		nextID: 1,
	}

	router.POST("/todos", store.createTodo)
	router.GET("/todos", store.listTodos)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.Run(":8080")
}

func (s *Store) createTodo(c *gin.Context) {
	var input struct {
		Title string `json:"title"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	s.Lock()
	defer s.Unlock()
	todo := Todo{
		ID:    s.nextID,
		Title: input.Title,
		Done:  false,
	}
	s.todos = append(s.todos, todo)
	s.nextID++
	c.JSON(http.StatusCreated, todo)
}

func (s *Store) listTodos(c *gin.Context) {
	s.Lock()
	defer s.Unlock()
	c.JSON(http.StatusOK, s.todos)
}
