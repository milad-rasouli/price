package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response[T any] struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    T      `json:"data"`
}

// Custom sends a custom response with the given status code, data and message
func Custom(c *gin.Context, statusCode int, data any, message string) {
	c.JSON(
		statusCode,
		Response[any]{
			Message: message,
			Status:  statusCode,
			Data:    data,
		},
	)
}

// Ok sends a success response with status 200
func Ok(c *gin.Context, data any, message string) {
	Custom(c, http.StatusOK, data, message)
}

// Created sends a success response with status 201
func Created(c *gin.Context, data any) {
	Custom(c, http.StatusCreated, data, "Created Successfully")
}

func NotFound(c *gin.Context) {
	Custom(c, http.StatusNotFound, nil, "not-found")
}

func InternalError(c *gin.Context) {
	Custom(c, http.StatusInternalServerError, nil, "internal-error")
}

func BadRequest(c *gin.Context, message string) {
	Custom(c, http.StatusBadRequest, nil, message)
}

func Pure(c *gin.Context, statusCode int, data any) {
	c.JSON(statusCode, data)
}
