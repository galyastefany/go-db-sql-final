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
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	// add
	store := NewParcelStore(db)
	parcel := getTestParcel()

	p, err := store.Add(parcel)
	parcel.Number = p
	require.NoError(t, err)
	require.NotEmpty(t, p)

	// get
	gp, err := store.Get(p)
	require.NoError(t, err)
	require.Equal(t, parcel, gp)

	// delete
	err = store.Delete(p)
	require.NoError(t, err)

	_, err = store.Get(p)
	require.ErrorIs(t, sql.ErrNoRows, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {

	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	// add
	store := NewParcelStore(db)
	parcel := getTestParcel()

	p, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, p)

	// set address
	newAddress := "new test address"
	err = store.SetAddress(p, newAddress)
	require.NoError(t, err)

	// check
	gp, err := store.Get(p)
	require.NoError(t, err)
	require.Equal(t, newAddress, gp.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {

	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	// add
	store := NewParcelStore(db)
	parcel := getTestParcel()
	p, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, p)

	// set status
	err = store.SetStatus(p, ParcelStatusSent)
	require.NoError(t, err)

	// check
	gp, err := store.Get(p)
	require.NoError(t, err)
	require.Equal(t, ParcelStatusSent, gp.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {

	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
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
		id, err := store.Add(parcels[i]) // добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
		require.NoError(t, err)
		require.NotEmpty(t, id)

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	require.Equal(t, len(parcels), len(storedParcels))

	// check
	for _, parcel := range storedParcels {

		v, ok := parcelMap[parcel.Number]
		require.True(t, ok)
		require.Equal(t, v, parcel)
	}
}
