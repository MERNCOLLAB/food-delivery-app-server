package address

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	http_helper "food-delivery-app-server/pkg/http"
)

type Handler struct {
	service *Service
}

func NewHandler(db *gorm.DB) *Handler {
	repo := NewRepository(db)
	service := NewService(repo)
	return &Handler{service: service}
}

func (h *Handler) CreateAddress(c *gin.Context) {
	req, err := http_helper.BindJSON[CreateAddressRequest](c)
	if err != nil {
		c.Error(err)
		return
	}

	userId, err := http_helper.ExtractUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	newAddr, err := h.service.CreateAddress(*req, userId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"message": "Address addded successfully", "address": newAddr})
}

func (h *Handler) GetAddress(c *gin.Context) {
	userId, err := http_helper.ExtractUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	addresses, err := h.service.GetAddress(userId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, addresses)
}

func (h *Handler) UpdateAddress(c *gin.Context) {
	addressId := c.Param("id")

	req, err := http_helper.BindJSON[UpdateAddressRequest](c)
	if err != nil {
		c.Error(err)
		return
	}

	userId, err := http_helper.ExtractUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	updatedAddr, err := h.service.UpdateAddress(addressId, userId, *req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"message": "Address updated successfully", "address": updatedAddr})
}

func (h *Handler) DeleteAddress(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Delete Address Endpoint"})
}
