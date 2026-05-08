package utils

import "github.com/gin-gonic/gin"

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func JSONResponse(c *gin.Context, statusCode int, status string, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}
