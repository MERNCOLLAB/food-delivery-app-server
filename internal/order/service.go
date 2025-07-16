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

func (s *Service) PlaceOrder(restaurantID string, userID string, orderReq PlaceOrderRequest) (*PlaceOrderResponse, error) {
	var totalAmount float64
	var orderItems []models.OrderItem

	restoID, err := utils.ParseId(restaurantID)
	if err != nil {
		return nil, appErr.NewBadRequest("Invalid Restaurant ID", err)
	}

	uID, err := utils.ParseId(userID)
	if err != nil {
		return nil, appErr.NewBadRequest("Invalid User ID", err)
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
		CustomerID:      &uID,
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

func (s *Service) GetOrderByRestaurant(restaurantID string) ([]models.Order, error) {
	restoID, err := utils.ParseId(restaurantID)
	if err != nil {
		return nil, appErr.NewBadRequest("Invalid restaurant ID", err)
	}

	orders, err := s.repo.GetOrdersByRestaurantID(restoID)
	if err != nil {
		return nil, appErr.NewInternal("Failed to get orders for restaurant", err)
	}

	return orders, nil
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

func (s *Service) GetOrderHistory(userId string) ([]models.Order, error) {
	uID, err := utils.ParseId(userId)
	if err != nil {
		return nil, appErr.NewBadRequest("Invalid user ID", err)
	}

	role, err := s.repo.GetUserRoleByID(uID)
	if err != nil {
		return nil, appErr.NewInternal("Failed to get user role", err)
	}

	orders, err := s.repo.GetOrderHistoryByUser(uID, role)
	if err != nil {
		return nil, appErr.NewInternal("Failed to get the order history", err)
	}

	return orders, nil
}

func (s *Service) GetAvailableOrders() {

}

func (s *Service) UpdateOrderStatus(req UpdateOrderStatusRequest, orderId string) error {
	orID, err := utils.ParseId(orderId)
	if err != nil {
		return appErr.NewBadRequest("Invalid order ID", err)
	}

	order, err := s.repo.GetOrderDetailsByID(orID)
	if err != nil {
		return appErr.NewInternal("Order not found", err)
	}

	currStatus := order.Status
	newStatus := models.Status(req.Status)

	allowedNext, ok := allowedStatusTransitions[currStatus]
	if !ok {
		return appErr.NewBadRequest("Invalid current order status", nil)
	}

	valid := false
	for _, s := range allowedNext {
		if s == newStatus {
			valid = true
			break
		}
	}

	if !valid {
		return appErr.NewBadRequest("Invalid status request", nil)
	}

	if err := s.repo.UpdateOrderStatus(orID, req.Status); err != nil {
		return appErr.NewInternal("Failed to update the order status", err)
	}

	return nil
}

func (s *Service) CancelOrder(orderId string, userId string) error {
	orID, err := utils.ParseId(orderId)
	if err != nil {
		return appErr.NewBadRequest("Invalid order ID", err)
	}

	uID, err := utils.ParseId(userId)
	if err != nil {
		return appErr.NewBadRequest("Invalid user ID", err)
	}

	order, err := s.repo.GetOrderDetailsByID(orID)
	if err != nil {
		return appErr.NewInternal("Order not found", err)
	}

	if order.Status != models.Pending && order.Status != models.AcceptedByOwner {
		return appErr.NewBadRequest("Order cannot be canceled at this stage", nil)
	}

	if order.CustomerID == nil || *order.CustomerID != uID {
		return appErr.NewUnauthorized("You are not allowed to cancel this order", nil)
	}

	if order.Status == models.AcceptedByOwner {
		if time.Since(order.UpdatedAt) > 3*time.Minute {
			return appErr.NewBadRequest("You can only cancel within 3 minutes after acceptance", nil)
		}
	}

	cancelStatus := string(models.Canceled)
	if err := s.repo.UpdateOrderStatus(orID, cancelStatus); err != nil {
		return appErr.NewInternal("Failed to cancel the order", err)
	}

	return nil
}
