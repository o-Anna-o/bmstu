package ds

import "time"

type RequestShip struct {
	ID                  int                 `gorm:"primaryKey;column:id"`
	Status              string              `gorm:"column:status"`
	CreatedAt           time.Time           `gorm:"column:created_at"`
	UserID              int                 `gorm:"column:user_id"`
	FormedAt            *time.Time          `gorm:"column:formed_at"`
	CompletedAt         *time.Time          `gorm:"column:completed_at"`
	Containers20ftCount int                 `gorm:"column:containers_20ft_count"`
	Containers40ftCount int                 `gorm:"column:containers_40ft_count"`
	Comment             string              `gorm:"column:comment"`
	LoadingSpeed        string              `gorm:"column:loading_speed"`
	LoadingTime         float64             `gorm:"column:loading_time"`
	Ships               []ShipInRequestShip `gorm:"foreignKey:RequestShipID"`
}

func (RequestShip) TableName() string {
	return "requests"
}
