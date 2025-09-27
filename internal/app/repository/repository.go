package repository

import (
	"fmt"
	"strings"
	"time"
)

type Repository struct {
	Ships    []Ship
	Requests map[int]Request
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
	ID         int     `gorm:"primaryKey"`
	Name       string  `gorm:"size:200;not null"`
	Capacity   float64 `gorm:"type:decimal(10,2)"`
	Length     float64 `gorm:"type:decimal(10,2)"`
	Width      float64 `gorm:"type:decimal(10,2)"`
	Draft      float64 `gorm:"type:decimal(10,2)"`
	Cranes     int
	Containers int
	Features   string `gorm:"type:text"`
	PhotoURL   string `gorm:"size:500"`
	IsActive   bool   `gorm:"default:true"`
}

type ShipInRequest struct {
	RequestID int       `gorm:"primaryKey"`
	ShipID    int       `gorm:"primaryKey"`
	Count     int       `gorm:"not null;default:1"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	Ship    Ship    `gorm:"foreignKey:ShipID"`
	Request Request `gorm:"foreignKey:RequestID"`
}

type Request struct {
	ID        int       `gorm:"primaryKey"`
	Status    string    `gorm:"type:varchar(20);default:'черновик'"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UserID    int       `gorm:"not null"`

	FormedAt    *time.Time `gorm:"null"`
	CompletedAt *time.Time `gorm:"null"`
	ModeratorID *int       `gorm:"null"`

	Containers20ftCount int     `gorm:"default:0"`
	Containers40ftCount int     `gorm:"default:0"`
	Comment             string  `gorm:"type:text"`
	LoadingSpeed        string  `gorm:"type:text"`
	LoadingTime         float64 `gorm:"type:decimal(10,2)"`

	Ships     []ShipInRequest `gorm:"foreignKey:RequestID"`
	User      User            `gorm:"foreignKey:UserID"`
	Moderator *User           `gorm:"foreignKey:ModeratorID"`
}

func NewRepository() (*Repository, error) {
	ships := []Ship{
		{
			ID:         1,
			Name:       "Ever Ace",
			Capacity:   23992,
			Length:     400,
			Width:      61.53,
			Draft:      17.0,
			Cranes:     6,
			Containers: 11996,
			Features:   "самый большой в мире, двигатель Wartsila 70950 кВт",
			PhotoURL:   "ever-ace.png",
		},
		{
			ID:         2,
			Name:       "FESCO Diomid",
			Capacity:   3108,
			Length:     195,
			Width:      32.20,
			Draft:      11.0,
			Cranes:     3,
			Containers: 536,
			Features:   "построен в 2010 г., судно класса Ice1 (для Арктики), дизельный двигатель, используется на Северном морском пути",
			PhotoURL:   "fesco-diomid.png",
		},
		{
			ID:         3,
			Name:       "HMM Algeciras",
			Capacity:   23964,
			Length:     399.9,
			Width:      61.0,
			Draft:      16.5,
			Cranes:     7,
			Containers: 11982,
			Features:   "двигатель MAN B&W 11G95ME-C9.5 мощностью 64 000 кВт, двойные двигатели, система рекуперации энергии, класс DNV GL",
			PhotoURL:   "hmm-algeciras.png",
		},
		{
			ID:         4,
			Name:       "MSC Gulsun",
			Capacity:   23756,
			Length:     399.9,
			Width:      61.4,
			Draft:      16.0,
			Cranes:     7,
			Containers: 11878,
			Features:   "первый в мире контейнеровоз, вмещающий более 23 000 TEU, двигатель MAN B&W 11G95ME-C9.5, класс DNV GL",
			PhotoURL:   "msc-gulsun.png",
		}}

	if len(ships) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	requests := make(map[int]Request)
	requests[1] = Request{ID: 1, Ships: []ShipInRequest{}}

	return &Repository{Ships: ships, Requests: requests}, nil
}

func (r *Repository) GetShips() ([]Ship, error) {
	// имитируем работу с БД
	ships := []Ship{
		{
			ID:         1,
			Name:       "Ever Ace",
			Capacity:   23992,
			Length:     400,
			Width:      61.53,
			Draft:      17.0,
			Cranes:     6,
			Containers: 11996,
			Features:   "самый большой в мире, двигатель Wartsila 70950 кВт",
			PhotoURL:   "ever-ace.png",
		},
		{
			ID:         2,
			Name:       "FESCO Diomid",
			Capacity:   3108,
			Length:     195,
			Width:      32.20,
			Draft:      11.0,
			Cranes:     3,
			Containers: 536,
			Features:   "построен в 2010 г., судно класса Ice1 (для Арктики), дизельный двигатель, используется на Северном морском пути",
			PhotoURL:   "fesco-diomid.png",
		},
		{
			ID:         3,
			Name:       "HMM Algeciras",
			Capacity:   23964,
			Length:     399.9,
			Width:      61.0,
			Draft:      16.5,
			Cranes:     7,
			Containers: 11982,
			Features:   "двигатель MAN B&W 11G95ME-C9.5 мощностью 64 000 кВт, двойные двигатели, система рекуперации энергии, класс DNV GL",
			PhotoURL:   "hmm-algeciras.png",
		},
		{
			ID:         4,
			Name:       "MSC Gulsun",
			Capacity:   23756,
			Length:     399.9,
			Width:      61.4,
			Draft:      16.0,
			Cranes:     7,
			Containers: 11878,
			Features:   "первый в мире контейнеровоз, вмещающий более 23 000 TEU, двигатель MAN B&W 11G95ME-C9.5, класс DNV GL",
			PhotoURL:   "msc-gulsun.png",
		},
	}

	if len(ships) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return r.Ships, nil
}

func (r *Repository) GetShip(id int) (Ship, error) {

	ships, err := r.GetShips()
	if err != nil {
		return Ship{}, err
	}

	for _, ship := range ships {
		if ship.ID == id {
			return ship, nil
		}
	}
	return Ship{}, fmt.Errorf("контейнеровоз не найден")
}

func (r *Repository) GetShipsByName(name string) ([]Ship, error) {
	ships, err := r.GetShips()
	if err != nil {
		return []Ship{}, err
	}

	var result []Ship
	for _, ship := range ships {
		if strings.Contains(strings.ToLower(ship.Name), strings.ToLower(name)) {
			result = append(result, ship)
		}
	}

	return result, nil
}

func (r *Repository) GetRequest(id int) (Request, error) {
	if request, ok := r.Requests[id]; ok {
		return request, nil
	}
	return Request{}, fmt.Errorf("заявка c id=%d не найдена", id)
}

func (r *Repository) RemoveShipFromRequest(requestID int, shipID int) error {
	request, ok := r.Requests[requestID]
	if !ok {
		return fmt.Errorf("заявка с id=%d не найдена", requestID)
	}
	for i, shipInRequest := range request.Ships {
		if shipInRequest.Ship.ID == shipID {
			request.Ships = append(request.Ships[:i], request.Ships[i+1:]...)
			r.Requests[requestID] = request
			return nil
		}
	}

	return fmt.Errorf("корабль с id=%d не найден в заявке", shipID)
}
