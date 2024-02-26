package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройка подключение к БД
	require.NoError(t, err)

	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавляем новую посылку в БД
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	// get
	// получаем только что добавленную посылку, проверяем отсутствие ошибки
	// проверяем, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
	res, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, err, parcel)

	// delete
	// удаляем добавленную посылку, проверяем отсутствие ошибки
	// проверяем, что посылку больше нельзя получить из БД
	err = store.Delete(id)
	require.NoError(t, err)
	res, err = store.Get(id)
	require.Empty(t, res)

}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройка подключения к БД
	require.NoError(t, err)

	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавляем новую посылку в БД, проверяем отсутствие ошибки
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	// set address
	// обновляем адрес, проверяем отсутствие ошибки
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	// check
	// получаем добавленную посылку
	res, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newAddress, res.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройка подключения к БД
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавляем новую посылку в БД, проверяем отсутствие ошибки
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	// set status
	// обновляем статус, проверяем отсутствие ошибки
	err = store.SetStatus(id, ParcelStatusSent)
	require.NoError(t, err)

	// check
	// получаем добавленную посылку
	res, err := store.Get(id)
	require.NoError(t, err)

	require.Equal(t, res.Status, ParcelStatusSent)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройка подключения к БД
	require.NoError(t, err)

	defer db.Close()

	store := NewParcelStore(db)
	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i]) // добавляем новую посылку в БД, проверяем отсутствие ошибки
		require.NoError(t, err)
		require.NotEmpty(t, id)

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client) // получаем список посылок по идентификатору клиента, сохранённого в переменной client
	require.NoError(t, err)
	require.Equal(t, len(parcels), len(storedParcels))
	// проверяем отсутствие ошибки
	// проверяем, что количество полученных посылок совпадает с количеством добавленных

	// check
	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// проверяем, что все посылки из storedParcels есть в parcelMap
		// проверяем, что значения полей полученных посылок заполнены верно
		value, ok := parcelMap[parcel.Number]
		require.True(t, ok)
		require.Equal(t, value, parcel)
	}
}
