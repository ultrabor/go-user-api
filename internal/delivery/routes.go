package delivery

import "github.com/gin-gonic/gin"

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/users", h.CreateUserHandler)
	r.GET("/users/:id", h.GetUserHandler)
	r.PUT("/users/:id", h.UpdateUserHandler)
	r.DELETE("/users/:id", h.DeleteUserHandler)
	r.GET("/users", h.GetAllUserHandler)
}
