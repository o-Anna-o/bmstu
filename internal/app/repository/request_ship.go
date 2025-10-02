package repository

import (
	"loading_time/internal/app/ds"

	"gorm.io/gorm"
)

func (r *Repository) GetRequestShip(id int) (ds.RequestShip, error) {
	request_ship := ds.RequestShip{}

	err := r.db.Preload("Ships.Ship").Where("id = ?", id).First(&request_ship).Error
	if err != nil {
		return ds.RequestShip{}, err
	}
	return request_ship, nil
}

func (r *Repository) GetOrCreateUserDraft(dummyUserID int) (ds.RequestShip, error) {
	var requestShip ds.RequestShip

	err := r.db.Preload("Ships.Ship").Where("status = ?", "черновик").First(&requestShip).Error
	if err == nil {
		return requestShip, nil
	}
	if err != gorm.ErrRecordNotFound {
		return ds.RequestShip{}, err
	}

	// Создаем новый черновик
	requestShip = ds.RequestShip{
		Status: "черновик",
	}

	err = r.db.Create(&requestShip).Error
	if err != nil {
		return ds.RequestShip{}, err
	}

	return requestShip, nil
}

// добавить корабль в заявку через ORM
func (r *Repository) AddShipToRequestShip(request_shipID, shipID int) error {
	var existingShip ds.ShipInRequestShip
	err := r.db.Where("request_id = ? AND ship_id = ?", request_shipID, shipID).First(&existingShip).Error

	if err == nil {
		existingShip.Count++
		return r.db.Save(&existingShip).Error
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	shipInRequestShip := ds.ShipInRequestShip{
		RequestShipID: request_shipID,
		ShipID:        shipID,
		Count:         1,
	}
	return r.db.Create(&shipInRequestShip).Error
}

// удалить корабль из заявки через ORM
func (r *Repository) RemoveShipFromRequestShip(request_shipID, shipID int) error {
	return r.db.Where("request_ship_id = ? AND ship_id = ?", request_shipID, shipID).Delete(&ds.ShipInRequestShip{}).Error
}

// удаление заявки через SQL
func (r *Repository) DeleteRequestShipSQL(request_shipID int) error {
	return r.db.Exec("UPDATE request_ships SET status = 'удалён' WHERE id = ?", request_shipID).Error
}

// вывестм заявку, исключая удаленные
func (r *Repository) GetRequestShipExcludingDeleted(id int) (ds.RequestShip, error) {
	var request_ship ds.RequestShip
	err := r.db.Preload("Ships.Ship").Where("id = ? AND status != ?", id, "удалён").First(&request_ship).Error
	if err != nil {
		return ds.RequestShip{}, err
	}
	return request_ship, nil
}
