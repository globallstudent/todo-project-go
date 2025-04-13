package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/globallstudent/todo-project-go/internal/models"
	"github.com/globallstudent/todo-project-go/internal/services"
)

type TodoHandler struct {
	todoService *services.TodoService
}

func NewTodoHandler(todoService *services.TodoService) *TodoHandler {
	return &TodoHandler{todoService: todoService}
}

// CreateTodo creates a new todo
// @Summary Create a new todo
// @Description Creates a new todo item for the authenticated user.
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param todo body models.Todo true "Todo data"
// @Success 201 {object} models.Todo
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Router /todos [post]

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var input models.Todo
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	input.UserID = userID.(int)

	if err := h.todoService.CreateTodo(c.Request.Context(), &input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, input)
}

func (h *TodoHandler) GetTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	todo, err := h.todoService.GetTodoByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	if role != "admin" && todo.UserID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// GetTodo retrieves a todo by ID
// @Summary Get a todo by ID
// @Description Retrieves a todo item by its ID. Users can only access their own todos unless they are admins.
// @Tags todos
// @Produce json
// @Security BearerAuth
// @Param id path int true "Todo ID"
// @Success 200 {object} models.Todo
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /todos/{id} [get]

func (h *TodoHandler) GetTodos(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var todos []*models.Todo
	var err error

	if role == "admin" {
		todos, err = h.todoService.GetTodosByUserID(c.Request.Context(), 0) // 0 = all users
	} else {
		todos, err = h.todoService.GetTodosByUserID(c.Request.Context(), userID.(int))
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todos)
}

// UpdateTodo updates a todo
// @Summary Update a todo
// @Description Updates an existing todo item. Users can only update their own todos unless they are admins.
// @Tags todos
// @Produce json
// @Security BearerAuth
// @Param id path int true "Todo ID"
// @Param todo body models.Todo true "Todo data"
// @Success 200 {object} models.Todo
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /todos/{id} [put]

func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input models.Todo
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.ID = id
	todo, err := h.todoService.GetTodoByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	if role != "admin" && todo.UserID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	if err := h.todoService.UpdateTodo(c.Request.Context(), &input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, input)
}

// DeleteTodo deletes a todo by ID
// @Summary Delete a todo by ID
// @Description Deletes a todo item by its ID. Users can only delete their own todos unless they are admins.
// @Tags todos
// @Produce json
// @Security BearerAuth
// @Param id path int true "Todo ID"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /todos/{id} [delete]

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	todo, err := h.todoService.GetTodoByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	if role != "admin" && todo.UserID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	if err := h.todoService.DeleteTodo(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "todo deleted"})
}
