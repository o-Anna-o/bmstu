package ds

type Ship struct {
	ID         int     `gorm:"primaryKey;column:id"`
	Name       string  `gorm:"column:name"`
	Capacity   float64 `gorm:"column:capacity"`
	Length     float64 `gorm:"column:length"`
	Width      float64 `gorm:"column:width"`
	Draft      float64 `gorm:"column:draft"`
	Cranes     int     `gorm:"column:cranes"`
	Containers int     `gorm:"column:containers"`
	Features   string  `gorm:"column:features"`
	PhotoURL   string  `gorm:"column:photo_url"`
	IsActive   bool    `gorm:"column:is_active"`
}

func (Ship) TableName() string {
	return "ships"
}
