package repository


import (
	"time" 
)


type NotificationLog struct {
	ID        uint      `gorm:"primarykey" json:"id"`                             
	EventType string    `gorm:"type:varchar(50);not null" json:"event_type"`      
	OrderID   string    `gorm:"type:varchar(100);not null;index" json:"order_id"` 
	Payload   string    `gorm:"type:text;not null" json:"payload"`                
	Message   string    `gorm:"type:text" json:"message"`                         
	CreatedAt time.Time `json:"created_at"`                                       
}
