package repository

import (
	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(product *Product) error
	GetProductByID(id uint) (*Product, error)
	GetProductByName(name string) (*Product, error)
	UpdateProduct(product *Product) error
	UpdateStock(productID uint, newStock int) error
	ListProducts(limit, offset int) ([]Product, error)
}

type StockReservationRepository interface {
	CreateReservation(reservation *StockReservation) error
	GetReservationByOrderID(orderID string) (*StockReservation, error)
	UpdateReservationStatus(reservationID uint, status ProductStatus) error
	DeleteReservation(reservationID uint) error
}
type productRepository struct {
	db *gorm.DB
}

type stockReservationRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}
func NewStockReservationRepository(db *gorm.DB) StockReservationRepository {
	return &stockReservationRepository{db: db}
}

func (r *productRepository) CreateProduct(product *Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) GetProductByID(id uint) (*Product, error) {
	var product Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetProductByName(name string) (*Product, error) {
	var product Product
	err := r.db.Where("name = ?", name).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}
func (r *productRepository) UpdateProduct(product *Product) error {
	return r.db.Save(product).Error
}
func (r *productRepository) UpdateStock(productID uint, newStock int) error {
	return r.db.Model(&Product{}).Where("id = ?", productID).Update("stock", newStock).Error
}
func (r *productRepository) ListProducts(limit, offset int) ([]Product, error) {
	var products []Product
	err := r.db.Limit(limit).Offset(offset).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
func (r *stockReservationRepository) CreateReservation(reservation *StockReservation) error {
	return r.db.Create(reservation).Error
}
func (r *stockReservationRepository) GetReservationByOrderID(orderID string) (*StockReservation, error) {
	var reservation StockReservation
	err := r.db.Where("order_id = ?", orderID).First(&reservation).Error
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}
func (r *stockReservationRepository) UpdateReservationStatus(reservationID uint, status ProductStatus) error {
	return r.db.Model(&StockReservation{}).Where("id = ?", reservationID).Update("status", status).Error
}
func (r *stockReservationRepository) DeleteReservation(reservationID uint) error {
	return r.db.Delete(&StockReservation{}, reservationID).Error
}
