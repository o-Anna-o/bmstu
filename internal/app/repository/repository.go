package repository

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type User struct {
	UserID              int     `gorm:"primaryKey;column:user_id"`
	FIO                 string  `gorm:"size:100;not null"`
	Login               string  `gorm:"size:100;unique;not null"`
	Password            string  `gorm:"size:10;not null"`
	Contacts            string  `gorm:"size:100"`
	CargoWeight         float64 `gorm:"type:decimal(10,2)"`
	Containers20ftCount int     `gorm:"default:0"`
	Containers40ftCount int     `gorm:"default:0"`
	IsModerator         bool    `gorm:"default:false"`
}

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
	IsActive   bool    `gorm:"column:is_active;default:true"`
}

type ShipInRequest struct {
	RequestID int       `gorm:"primaryKey;column:request_id"`
	ShipID    int       `gorm:"primaryKey;column:ship_id"`
	Count     int       `gorm:"column:count;default:1"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`

	Ship Ship `gorm:"foreignKey:ShipID"`
}

type Request struct {
	ID        int       `gorm:"primaryKey;column:id"`
	Status    string    `gorm:"column:status;default:'черновик'"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UserID    int       `gorm:"column:user_id"`

	FormedAt    *time.Time `gorm:"column:formed_at"`
	CompletedAt *time.Time `gorm:"column:completed_at"`
	ModeratorID *int       `gorm:"column:moderator_id"`

	Containers20ftCount int     `gorm:"column:containers_20ft_count;default:0"`
	Containers40ftCount int     `gorm:"column:containers_40ft_count;default:0"`
	Comment             string  `gorm:"column:comment"`
	LoadingSpeed        string  `gorm:"column:loading_speed"`
	LoadingTime         float64 `gorm:"column:loading_time"`

	Ships []ShipInRequest `gorm:"foreignKey:RequestID"`
}

func NewRepository(db *gorm.DB) (*Repository, error) {
	return &Repository{db: db}, nil
}

func (r *Repository) GetShips() ([]Ship, error) {
	var ships []Ship
	result := r.db.Where("is_active = ?", true).Find(&ships)
	if result.Error != nil {
		return nil, result.Error
	}
	return ships, nil
}

func (r *Repository) GetShip(id int) (Ship, error) {
	var ship Ship
	result := r.db.Where("id = ? AND is_active = ?", id, true).First(&ship)
	if result.Error != nil {
		return Ship{}, result.Error
	}
	return ship, nil
}

func (r *Repository) GetShipsByName(name string) ([]Ship, error) {
	var ships []Ship
	result := r.db.Where("name ILIKE ? AND is_active = ?", "%"+name+"%", true).Find(&ships)
	if result.Error != nil {
		return nil, result.Error
	}
	return ships, nil
}

func (r *Repository) GetRequest(id int) (Request, error) {
	var request Request
	result := r.db.Preload("Ships.Ship").Where("id = ? AND status != 'удалена'", id).First(&request)
	if result.Error != nil {
		return Request{}, result.Error
	}
	return request, nil
}

func (r *Repository) RemoveShipFromRequest(requestID int, shipID int) error {
	result := r.db.Exec("DELETE FROM request_ships WHERE request_id = ? AND ship_id = ?", requestID, shipID)
	return result.Error
}

func (r *Repository) DeleteRequest(requestID int) error {
	result := r.db.Exec("UPDATE requests SET status = 'удалена' WHERE id = ?", requestID)
	return result.Error
}

func (r *Repository) AddToRequest(requestID int, shipID int) error {
	// Проверяем, существует ли уже такой корабль в заявке
	var existingShip ShipInRequest
	result := r.db.Where("request_id = ? AND ship_id = ?", requestID, shipID).First(&existingShip)

	if result.Error == nil {
		// Обновляем количество
		return r.db.Exec("UPDATE request_ships SET count = count + 1 WHERE request_id = ? AND ship_id = ?",
			requestID, shipID).Error
	}

	// Создаем новую запись
	return r.db.Exec("INSERT INTO request_ships (request_id, ship_id, count, created_at) VALUES (?, ?, 1, ?)",
		requestID, shipID, time.Now()).Error
}

// НОВЫЙ МЕТОД: Получение или создание черновой заявки
func (r *Repository) GetOrCreateDraftRequest(userID int) (Request, error) {
	var request Request

	// Пытаемся найти существующую черновую заявку
	result := r.db.Preload("Ships.Ship").
		Where("user_id = ? AND status = 'черновик'", userID).
		First(&request)

	if result.Error == nil {
		return request, nil // Заявка найдена
	}

	// Создаем новую заявку если не найдена
	request = Request{
		Status:    "черновик",
		UserID:    userID,
		Comment:   "Демо-заявка",
		CreatedAt: time.Now(),
	}

	result = r.db.Create(&request)
	return request, result.Error
}

// НОВЫЙ МЕТОД: Подсчет количества кораблей в заявке
func (r *Repository) GetRequestCount(userID int) (int, error) {
	var request Request
	result := r.db.Preload("Ships").
		Where("user_id = ? AND status = 'черновик'", userID).
		First(&request)

	if result.Error != nil {
		// Если заявки нет, возвращаем 0
		return 0, nil
	}

	// Считаем общее количество кораблей в заявке
	count := 0
	for _, ship := range request.Ships {
		count += ship.Count
	}

	return count, nil
}

// ______________________________________________________________________
func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Автомиграция
	err = db.AutoMigrate(&Ship{}, &Request{}, &ShipInRequest{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}
