package delivery

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ultrabor/go-user-api/internal/domain"
	m "github.com/ultrabor/go-user-api/internal/domain"
)

type UserHandler struct {
	service domain.UserService
}

func NewUserHandler(service domain.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (handler *UserHandler) CreateUserHandler(c *gin.Context) {

	var input struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(err)
		return
	}

	user, err := handler.service.CreateUser(input.Name, input.Age)

	if err != nil {
		c.Error(err)
	}

	c.JSON(http.StatusOK, gin.H{"id": user.ID})
}

func (handler *UserHandler) DeleteUserHandler(c *gin.Context) {

	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.Error(err)
		return
	}

	err = handler.service.DeleteUser(id)

	if err != nil {
		c.Error(err)
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (handler *UserHandler) GetUserHandler(c *gin.Context) {

	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.Error(err)
		return
	}

	user, err := handler.service.GetUser(id)
	if err != nil {
		c.Error(err)
	}

	c.JSON(http.StatusFound, user)

}

func (h *UserHandler) GetAllUserHandler(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	pageStr := c.DefaultQuery("page", "1")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 200 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		return
	}

	name := c.Query("name")

	var age int
	if ageStr := c.Query("age"); ageStr != "" {
		age, err = strconv.Atoi(ageStr)
		if err != nil || age < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid age"})
			return
		}
	}

	users, err := h.service.GetAllUsers(limit, page, name, age)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (handler *UserHandler) UpdateUserHandler(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.Error(err)
		return
	}

	var input struct {
		Name *string `json:"name"`
		Age  *int    `json:"age"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(err)
		return
	}

	user := m.User{ID: id, Name: "", Age: 0}

	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Age != nil {
		user.Age = *input.Age
	}

	updatedUser, err := handler.service.UpdateUser(user)

	if err != nil {
		c.Error(err)
	}

	c.JSON(http.StatusOK, updatedUser)

}
