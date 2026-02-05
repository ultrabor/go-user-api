package delivery

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ultrabor/go-user-api/internal/domain"
	"github.com/ultrabor/go-user-api/internal/services"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func parseID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, false
	}
	return id, true
}

func (h *UserHandler) CreateUserHandler(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
		Age  int    `json:"age"  binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Age <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "age must be > 0"})
		return
	}

	user, err := h.service.CreateUser(input.Name, input.Age)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": user.ID})
}

func (h *UserHandler) DeleteUserHandler(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	if err := h.service.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) GetUserHandler(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	user, err := h.service.GetUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetAllUserHandler(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	pageStr := c.DefaultQuery("page", "1")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	name := c.Query("name")

	var age int
	if ageStr := c.Query("age"); ageStr != "" {
		age, err = strconv.Atoi(ageStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
	}

	users, err := h.service.GetAllUsers(limit, page, name, age)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) UpdateUserHandler(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var input struct {
		Name *string `json:"name"`
		Age  *int    `json:"age"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Age != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "age must be > 0"})
		return
	}
	if input.Name == nil && input.Age == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	existing, err := h.service.GetUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	if input.Name != nil {
		existing.Name = *input.Name
	}
	if input.Age != nil {
		existing.Age = *input.Age
	}

	updated, err := h.service.UpdateUser(existing)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	c.JSON(http.StatusOK, updated)
}
