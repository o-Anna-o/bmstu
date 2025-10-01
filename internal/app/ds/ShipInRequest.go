package ds

import "time"

type ShipInRequestShip struct {
	RequestShipID int       `gorm:"primaryKey;column:request_ship_id"`
	ShipID    int       `gorm:"primaryKey;column:ship_id"`
	Count     int       `gorm:"column:count"`
	CreatedAt time.Time `gorm:"column:created_at"`
	Ship      Ship      `gorm:"foreignKey:ShipID"`
}

func (ShipInRequestShip) TableName() string {
	return "request_ship_ships"
}
