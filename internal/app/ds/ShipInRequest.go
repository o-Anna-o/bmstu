package ds

import "time"

type ShipInRequest struct {
	RequestID int       `gorm:"primaryKey;column:request_id"`
	ShipID    int       `gorm:"primaryKey;column:ship_id"`
	Count     int       `gorm:"column:count"`
	CreatedAt time.Time `gorm:"column:created_at"`
	Ship      Ship      `gorm:"foreignKey:ShipID"`
}

func (ShipInRequest) TableName() string {
	return "request_ships"
}
