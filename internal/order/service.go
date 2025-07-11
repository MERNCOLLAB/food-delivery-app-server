package order

import (
	appErr "food-delivery-app-server/pkg/errors"
	"time"

	"food-delivery-app-server/models"
	"food-delivery-app-server/pkg/utils"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) PlaceOrder(restaurantID string, orderReq PlaceOrderRequest) (*PlaceOrderResponse, error) {
	var totalAmount float64
	var orderItems []models.OrderItem

	restoID, err := utils.ParseId(restaurantID)
	if err != nil {
		return nil, appErr.NewBadRequest("Invalid Restaurant ID", err)
	}

	for _, item := range orderReq.Items {
		menuItem, err := s.repo.GetMenuItemDetails(item.MenuItem)
		if err != nil {
			return nil, err
		}
		itemTotal := menuItem.Price * float64(item.Quantity)
		totalAmount += itemTotal

		orderItems = append(orderItems, models.OrderItem{
			ID:         utils.GenerateUUID(),
			MenuItemID: menuItem.ID,
			Quantity:   int32(item.Quantity),
			Price:      menuItem.Price,
		})
	}

	deliveryFee := 50.0
	totalAmount += deliveryFee

	orderID := utils.GenerateUUID()

	order := &models.Order{
		ID:              orderID,
		RestaurantID:    restoID,
		Status:          models.Status("PENDING"),
		TotalAmount:     totalAmount,
		DeliveryFee:     deliveryFee,
		DeliveryAddress: orderReq.DeliveryAddress,
		PlacedAt:        time.Now(),
	}

	createdOrder, err := s.repo.CreateOrder(order, orderItems)
	if err != nil {
		return nil, err
	}

	resto, err := s.repo.GetRestaurantByID(restoID)
	if err == nil && resto.OwnerID != uuid.Nil {
		notification := &models.Notification{
			ID:        utils.GenerateUUID(),
			UserID:    resto.OwnerID,
			OrderID:   &createdOrder.ID,
			Message:   "You have a new order!",
			IsRead:    false,
			CreatedAt: time.Now(),
		}

		_ = s.repo.CreateNotification(notification)
	}

	orderRes := &PlaceOrderResponse{
		OrderID:         createdOrder.ID,
		Status:          createdOrder.Status,
		TotalAmount:     createdOrder.TotalAmount,
		DeliveryFee:     createdOrder.DeliveryFee,
		DeliveryAddress: createdOrder.DeliveryAddress,
		PlacedAt:        createdOrder.PlacedAt.Format(time.RFC1123),
		Items:           orderItems,
	}

	return orderRes, nil
}

func (s *Service) GetOrderByRestaurant() {

}

func (s *Service) GetOrderDetails(orderId string) (*models.Order, error) {
	orID, err := utils.ParseId(orderId)
	if err != nil {
		return nil, appErr.NewInternal("Invalid ID", err)
	}

	order, err := s.repo.GetOrderDetailsByID(orID)
	if err != nil {
		return nil, appErr.NewInternal("Failed to query order details", err)
	}

	return order, nil
}

func (s *Service) GetOrderHistory() {

}

func (s *Service) GetAvailableOrders() {

}
