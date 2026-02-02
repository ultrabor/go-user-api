package delivery

import "github.com/gin-gonic/gin"

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/products", h.CreateUserHandler)
	r.GET("/products/:id", h.GetUserHandler)
	r.PUT("/products/:id", h.UpdateUserHandler)
	r.DELETE("/products/:id", h.DeleteUserHandler)
	r.GET("/products", h.GetAllUserHandler)
}
