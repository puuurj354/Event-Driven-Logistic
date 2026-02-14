package repository


import (
	"gorm.io/gorm" 
)

type NotificationRepository interface {
	SaveLog(log *NotificationLog) error                         
	GetLogsByOrderID(orderID string) ([]NotificationLog, error) 
	GetRecentLogs(limit int) ([]NotificationLog, error)         
}


type notificationRepository struct {
	db *gorm.DB 
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db} 
}
func (r *notificationRepository) SaveLog(log *NotificationLog) error {
	result := r.db.Create(log) 
	return result.Error        
}


func (r *notificationRepository) GetLogsByOrderID(orderID string) ([]NotificationLog, error) {
	var logs []NotificationLog                                                         
	result := r.db.Where("order_id = ?", orderID).Order("created_at DESC").Find(&logs) 
	return logs, result.Error                                                          
}


func (r *notificationRepository) GetRecentLogs(limit int) ([]NotificationLog, error) {
	var logs []NotificationLog                                       
	result := r.db.Order("created_at DESC").Limit(limit).Find(&logs) 
	return logs, result.Error                                        
}
