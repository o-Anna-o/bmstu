package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Ship struct { // вот наша новая структура
	ID    int    // поля структур, которые передаются в шаблон
	Title string // ОБЯЗАТЕЛЬНО должны быть написаны с заглавной буквы (то есть публичными)
}

func (r *Repository) GetShips() ([]Ship, error) {
	// имитируем работу с БД. Типа мы выполнили sql запрос и получили эти строки из БД
	ships := []Ship{ // массив элементов из наших структур
		{
			ID:    1,
			Title: "first ship",
		},
		{
			ID:    2,
			Title: "second ship",
		},
		{
			ID:    3,
			Title: "third ship",
		},
	}
	// обязательно проверяем ошибки, и если они появились - передаем выше, то есть хендлеру
	// тут я снова искусственно обработаю "ошибку" чисто чтобы показать вам как их передавать выше
	if len(ships) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return ships, nil
}

func (r *Repository) GetShip(id int) (Ship, error) {
	// тут у вас будет логика получения нужной услуги, тоже наверное через цикл в первой лабе, и через запрос к БД начиная со второй
	ships, err := r.GetShips()
	if err != nil {
		return Ship{}, err // тут у нас уже есть кастомная ошибка из нашего метода, поэтому мы можем просто вернуть ее
	}

	for _, ship := range ships {
		if ship.ID == id {
			return ship, nil // если нашли, то просто возвращаем найденный заказ (услугу) без ошибок
		}
	}
	return Ship{}, fmt.Errorf("заказ не найден") // тут нужна кастомная ошибка, чтобы понимать на каком этапе возникла ошибка и что произошло
}

func (r *Repository) GetShipsByTitle(title string) ([]Ship, error) {
	ships, err := r.GetShips()
	if err != nil {
		return []Ship{}, err
	}

	var result []Ship
	for _, ship := range ships {
		if strings.Contains(strings.ToLower(ship.Title), strings.ToLower(title)) {
			result = append(result, ship)
		}
	}

	return result, nil
}
