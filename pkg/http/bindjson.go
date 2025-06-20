package http_helper

import (
	appErr "food-delivery-app-server/pkg/errors"

	"github.com/gin-gonic/gin"
)

func BindJSON[T any](c *gin.Context) (*T, error) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, appErr.NewBadRequest("Invalid JSON request body", err)
	}

	return &req, nil
}
