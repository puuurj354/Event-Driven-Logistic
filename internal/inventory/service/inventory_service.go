package service

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/purnama/Event-Driven-Logistic/internal/inventory/repository"
)

// InventoryService mendefinisikan kontrak bisnis logic untuk Inventory.
// Semua method mengembalikan error jika operasi gagal.
type InventoryService interface {
	// CreateProduct menambahkan produk baru ke inventory.
	CreateProduct(name string, stock int) (*repository.Product, error)

	// ListProducts mengambil daftar produk dengan pagination.
	ListProducts(limit, offset int) ([]repository.Product, error)

	// GetProductByID mengambil produk berdasarkan ID.
	GetProductByID(id uint) (*repository.Product, error)

	// ReserveStock mereservasi stok untuk sebuah order.
	ReserveStock(orderID uuid.UUID, productName string, quantity int) (*repository.StockReservation, error)

	// ConfirmReservation mengkonfirmasi reservasi stok.
	ConfirmReservation(orderID string) error

	// ReleaseReservation mengembalikan stok yang direservasi.
	ReleaseReservation(orderID string) error
}

type inventoryService struct {
	productRepo     repository.ProductRepository
	reservationRepo repository.StockReservationRepository
}

func NewInventoryService(
	productRepo repository.ProductRepository,
	reservationRepo repository.StockReservationRepository,
) InventoryService {
	return &inventoryService{
		productRepo:     productRepo,
		reservationRepo: reservationRepo,
	}
}

// CreateProduct menambahkan produk baru ke inventory.
// Parameter: name — nama produk (wajib), stock — jumlah stok awal (>= 0).
// Return: product yang berhasil dibuat, atau error jika validasi gagal.
func (s *inventoryService) CreateProduct(name string, stock int) (*repository.Product, error) {
	// Validasi nama produk — tidak boleh kosong
	if name == "" {
		return nil, errors.New("nama produk tidak boleh kosong") // Guard clause
	}

	// Validasi stok awal — tidak boleh negatif
	if stock < 0 {
		return nil, errors.New("stok awal tidak boleh negatif") // Guard clause
	}

	// Buat entity produk baru
	product := &repository.Product{
		Name:  name,  // Nama produk dari parameter
		Stock: stock, // Stok awal dari parameter
	}

	// Simpan ke database via repository
	if err := s.productRepo.CreateProduct(product); err != nil {
		log.Printf("❌ Gagal membuat produk: name=%s, error=%v", name, err) // Log error
		return nil, fmt.Errorf("gagal membuat produk: %w", err)            // Wrap error
	}

	log.Printf("✅ Produk berhasil dibuat: ID=%d, Name=%s, Stock=%d", product.ID, name, stock) // Log sukses
	return product, nil                                                                       // Return produk baru
}

func (s *inventoryService) ListProducts(limit, offset int) ([]repository.Product, error) {

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	products, err := s.productRepo.ListProducts(limit, offset)
	if err != nil {
		log.Printf("❌ Gagal mengambil daftar produk: error=%v", err)
		return nil, err
	}
	return products, nil
}

func (s *inventoryService) GetProductByID(id uint) (*repository.Product, error) {
	product, err := s.productRepo.GetProductByID(id)
	if err != nil {
		log.Printf("❌ Produk tidak ditemukan: ID=%d, error=%v", id, err)
		return nil, err
	}
	return product, nil
}

func (s *inventoryService) ReserveStock(orderID uuid.UUID, productName string, quantity int) (*repository.StockReservation, error) {

	if quantity <= 0 {
		return nil, errors.New("quantity harus lebih dari 0")
	}

	product, err := s.productRepo.GetProductByName(productName)
	if err != nil {
		return nil, fmt.Errorf("produk '%s' tidak ditemukan: %w", productName, err)
	}

	if product.Stock < quantity {
		return nil, fmt.Errorf("stok tidak cukup untuk '%s': tersedia=%d, diminta=%d",
			productName, product.Stock, quantity)
	}

	newStock := product.Stock - quantity
	if err := s.productRepo.UpdateStock(product.ID, newStock); err != nil {
		return nil, fmt.Errorf("gagal mengurangi stok: %w", err)
	}

	reservation := &repository.StockReservation{
		OrderID:   orderID,
		ProductID: product.ID,
		Quantity:  quantity,
		Status:    repository.ProductStatusReserved,
	}

	if err := s.reservationRepo.CreateReservation(reservation); err != nil {

		_ = s.productRepo.UpdateStock(product.ID, product.Stock)
		return nil, fmt.Errorf("gagal membuat reservasi stok: %w", err)
	}

	log.Printf("✅ Stok direservasi: OrderID=%s, Product=%s, Qty=%d, Sisa=%d",
		orderID, productName, quantity, newStock)

	return reservation, nil
}

func (s *inventoryService) ConfirmReservation(orderID string) error {

	reservation, err := s.reservationRepo.GetReservationByOrderID(orderID)
	if err != nil {
		return fmt.Errorf("reservasi tidak ditemukan untuk order: %s, error: %w", orderID, err)
	}

	if err := s.reservationRepo.UpdateReservationStatus(reservation.ID, repository.ProductStatusConfirmed); err != nil {
		return fmt.Errorf("gagal konfirmasi reservasi: %w", err)
	}

	log.Printf("✅ Reservasi dikonfirmasi: OrderID=%s, ReservationID=%d", orderID, reservation.ID)
	return nil
}
func (s *inventoryService) ReleaseReservation(orderID string) error {

	reservation, err := s.reservationRepo.GetReservationByOrderID(orderID)
	if err != nil {
		return fmt.Errorf("reservasi tidak ditemukan untuk order: %s, error: %w", orderID, err)
	}

	product, err := s.productRepo.GetProductByID(reservation.ProductID)
	if err != nil {
		return fmt.Errorf("produk tidak ditemukan: ID=%d, error: %w", reservation.ProductID, err)
	}

	newStock := product.Stock + reservation.Quantity
	if err := s.productRepo.UpdateStock(product.ID, newStock); err != nil {
		return fmt.Errorf("gagal mengembalikan stok: %w", err)
	}

	if err := s.reservationRepo.UpdateReservationStatus(reservation.ID, repository.ProductStatusReleased); err != nil {
		return fmt.Errorf("gagal update status reservasi: %w", err)
	}

	log.Printf("✅ Stok dikembalikan: OrderID=%s, Product=%s, Qty=%d, StokBaru=%d",
		orderID, product.Name, reservation.Quantity, newStock)

	return nil
}
