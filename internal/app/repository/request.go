package repository

import (
	"loading_time/internal/app/ds"
)

func (r *Repository) GetRequest(id int) (ds.Request, error) {
	request := ds.Request{}
	// обязательно проверяем ошибки, и если они появились - передаем выше, то есть хендлеру
	err := r.db.Preload("Ships.Ship").Where("id = ?", id).First(&request).Error
	if err != nil {
		return ds.Request{}, err
	}
	return request, nil
}

// GetOrCreateUserDraft - получить или создать черновую заявку пользователя
func (r *Repository) GetOrCreateUserDraft(userID int) (ds.Request, error) {
	var request ds.Request

	// Ищем существующую черновую заявку С ПОДГРУЗКОЙ КОРАБЛЕЙ
	err := r.db.Preload("Ships.Ship").Where("user_id = ? AND status = ?", userID, "черновик").First(&request).Error
	if err == nil {
		return request, nil // Заявка найдена
	}

	// Создаем новую черновую заявку
	request = ds.Request{
		Status: "черновик",
		UserID: userID,
	}

	err = r.db.Create(&request).Error
	if err != nil {
		return ds.Request{}, err
	}

	// После создания снова загружаем с подгрузкой кораблей
	err = r.db.Preload("Ships.Ship").Where("id = ?", request.ID).First(&request).Error
	if err != nil {
		return ds.Request{}, err
	}

	return request, nil
}

// AddShipToRequest - добавить корабль в заявку через ORM
func (r *Repository) AddShipToRequest(requestID, shipID int) error {
	// Сначала проверяем, есть ли уже такой корабль в заявке
	var existingShip ds.ShipInRequest
	err := r.db.Where("request_id = ? AND ship_id = ?", requestID, shipID).First(&existingShip).Error

	if err == nil {
		// Корабль уже есть в заявке - увеличиваем количество
		existingShip.Count++
		return r.db.Save(&existingShip).Error
	}

	// Корабля нет в заявке - создаем новую запись
	shipInRequest := ds.ShipInRequest{
		RequestID: requestID,
		ShipID:    shipID,
		Count:     1,
	}
	return r.db.Create(&shipInRequest).Error
}

// RemoveShipFromRequest - удалить корабль из заявки
func (r *Repository) RemoveShipFromRequest(requestID, shipID int) error {
	return r.db.Where("request_id = ? AND ship_id = ?", requestID, shipID).Delete(&ds.ShipInRequest{}).Error
}
