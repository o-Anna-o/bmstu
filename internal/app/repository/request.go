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
