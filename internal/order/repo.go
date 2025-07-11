package order

import (
	"food-delivery-app-server/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateOrder(order *models.Order, items []models.OrderItem) (*models.Order, error) {
	// Starts a transaction to ensure order and order items creation are safe. If any one of those fails, the entire transaction
	// rolls back to maintain data integrity
	tx := r.db.Begin()
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for i := range items {
		items[i].OrderID = order.ID
		if err := tx.Create(&items[i]).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()
	return order, nil
}

func (r *Repository) GetMenuItemDetails(menuItemID uuid.UUID) (*models.MenuItem, error) {
	var menuItem models.MenuItem
	if err := r.db.Where("id = ?", menuItemID).First(&menuItem).Error; err != nil {
		return nil, err
	}
	return &menuItem, nil
}

func (r *Repository) GetRestaurantByID(restoId uuid.UUID) (*models.Restaurant, error) {
	var resto models.Restaurant
	if err := r.db.Where("id = ?", restoId).First(&resto).Error; err != nil {
		return nil, err
	}
	return &resto, nil
}

func (r *Repository) CreateNotification(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

func (r *Repository) GetOrderDetailsByID(orID uuid.UUID) (*models.Order, error) {
	var order models.Order

	err := r.db.
		Preload("OrderItems.MenuItem").
		Preload("Restaurant").
		Preload("Customer").
		Preload("Driver").
		First(&order, "id = ?", orID).Error

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *Repository) GetOrderByRestaurantID() {

}

func (r *Repository) UpdateOrderStatus() {

}
