package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	res, errAdd := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)", p.Client, p.Status, p.Address, p.CreatedAt)
	if errAdd != nil {
		return 0, errAdd
	}
	number, errNumber := res.LastInsertId()
	if errNumber != nil {
		return 0, errNumber
	}

	// верните идентификатор последней добавленной записи
	return int(number), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = ?", number)

	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	errAdd := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if errAdd != nil {
		return p, errAdd
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, errGetByClient := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = ?", client)
	if errGetByClient != nil {
		return nil, errGetByClient
	}
	defer rows.Close()

	// заполните срез Parcel данными из таблицы
	var res []Parcel
	p := Parcel{}
	for rows.Next() {
		errScanRows := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if errScanRows != nil {
			return res, errScanRows
		}
		res = append(res, p)
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, errSetStatus := s.db.Exec("UPDATE parcel SET status = ? WHERE number = ?", status, number)
	if errSetStatus != nil {
		return errSetStatus
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	_, errSetAddress := s.db.Exec("UPDATE parcel SET address = ? WHERE number = ? AND status = ?", address, number, ParcelStatusRegistered)
	if errSetAddress != nil {
		return errSetAddress
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	_, errDelete := s.db.Exec("DELETE FROM parcel WHERE number = ? AND status = ?", number, ParcelStatusRegistered)
	if errDelete != nil {
		return errDelete
	}

	return nil
}
