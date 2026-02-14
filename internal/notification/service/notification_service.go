package service

import (
	"fmt"
	"log"

	"github.com/purnama/Event-Driven-Logistic/internal/notification/repository"
)

type Broadcaster interface {
	Broadcast(message []byte)
}

type NotificationService interface {
	ProcessEvent(eventType, orderID, payload string) error
	GetRecentLogs(limit int) ([]repository.NotificationLog, error)
	GetLogsByOrderID(orderID string) ([]repository.NotificationLog, error)
}

type notificationService struct {
	repo repository.NotificationRepository
	hub  Broadcaster
}

func NewNotificationService(repo repository.NotificationRepository, hub Broadcaster) NotificationService {
	return &notificationService{
		repo: repo,
		hub:  hub,
	}
}

func humanMessage(eventType string) string {
	switch eventType {
	case "order.created":
		return "📦 Pesanan baru dibuat"
	case "payment.success":
		return "💳 Pembayaran berhasil"
	case "payment.failed":
		return "❌ Pembayaran gagal"
	case "stock.reserved":
		return "📦 Stok berhasil direservasi"
	case "stock.failed":
		return "⚠️ Stok tidak tersedia"
	default:
		return fmt.Sprintf("🔔 Event: %s", eventType)
	}
}

func (s *notificationService) ProcessEvent(eventType, orderID, payload string) error {
	notifLog := &repository.NotificationLog{
		EventType: eventType,
		OrderID:   orderID,
		Payload:   payload,
		Message:   humanMessage(eventType),
	}

	if err := s.repo.SaveLog(notifLog); err != nil {
		log.Printf("❌ Gagal menyimpan notification log: %v", err)
		return err
	}
	log.Printf("✅ Notification log disimpan: [%s] order=%s", eventType, orderID)

	wsMessage := fmt.Sprintf(`{"event_type":"%s","order_id":"%s","message":"%s","payload":%s}`,
		eventType,
		orderID,
		humanMessage(eventType),
		payload,
	)
	s.hub.Broadcast([]byte(wsMessage))

	return nil
}
func (s *notificationService) GetRecentLogs(limit int) ([]repository.NotificationLog, error) {
	return s.repo.GetRecentLogs(limit)
}

func (s *notificationService) GetLogsByOrderID(orderID string) ([]repository.NotificationLog, error) {
	return s.repo.GetLogsByOrderID(orderID)
}
