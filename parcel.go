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
	// добавление строки в таблицу
	res, err := s.db.Exec("insert into parcel (client, status, address, created_at) values (:client,:status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))

	if err != nil {
		return 0, err
	}

	// идентификатор последней добавленной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// чтение строки по заданному number
	// возвращается только одна строка
	// заполняем Parcel данными из таблицы
	p := Parcel{}
	row := s.db.QueryRow("select number, client, status, address, created_at from Parcel where number =:number", sql.Named("number", number))
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// чтение строки по заданному client
	// возвращается несколько строк
	// заполняем Parcel данными из таблицы
	var res []Parcel

	rows, err := s.db.Query(
		"SELECT number, client, status, address, created_at FROM parcel WHERE client = :client",
		sql.Named("client", client))

	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		p := Parcel{}

		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return res, err
		}
		res = append(res, p)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуем обновление статуса в таблице parcel
	_, err := s.db.Exec("update parcel set status =:status where number=:number",
		sql.Named("status", status), sql.Named("number", number))
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуем обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	_, err := s.db.Exec("update parcel set address = :address where number =:number and status=:status",
		sql.Named("address", address), sql.Named("number", number), sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуем удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	_, err := s.db.Exec("delete from parcel where number = :number and status = :status",
		sql.Named("number", number), sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return err
	}
	return nil
}
