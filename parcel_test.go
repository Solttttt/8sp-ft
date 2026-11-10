package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	store := NewParcelStore(db)
	parcel := getTestParcel()
	defer db.Close()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки
	number, err := store.Add(parcel)
	require.NoError(t, err)
	//и наличии идентификатора
	assert.NotEqual(t, 0, number)

	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	p, err := store.Get(number)
	require.NoError(t, err)
	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
	assert.Equal(t, number, p.Number)
	assert.Equal(t, parcel.Client, p.Client)
	assert.Equal(t, parcel.Status, p.Status)
	assert.Equal(t, parcel.Address, p.Address)
	assert.Equal(t, parcel.CreatedAt, p.CreatedAt)

	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	err = store.Delete(number)
	require.NoError(t, err)
	// проверьте, что посылку больше нельзя получить из БД
	_, err = store.Get(number)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	store := NewParcelStore(db)
	parcel := getTestParcel()
	defer db.Close()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки
	number, err := store.Add(parcel)
	require.NoError(t, err)
	//и наличии идентификатора
	assert.NotEqual(t, 0, number)

	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"
	err = store.SetAddress(number, newAddress)
	require.NoError(t, err)
	// check
	// получите добавленную посылку и убедитесь, что адрес обновился
	p, err := store.Get(number)
	require.NoError(t, err)
	assert.Equal(t, newAddress, p.Address)

	// set address для Sent посылки
	// установите статус sent и попробуйте снова обновить адрес
	err = store.SetStatus(number, ParcelStatusSent)
	require.NoError(t, err)
	newAddress2 := "new test2 address"
	err = store.SetAddress(number, newAddress2)
	require.NoError(t, err)
	// check
	p, err = store.Get(number)
	require.NoError(t, err)
	assert.Equal(t, newAddress, p.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	store := NewParcelStore(db)
	parcel := getTestParcel()
	defer db.Close()

	// add
	// добавьте новую 	посылку в БД, убедитесь в отсутствии ошибки
	number, err := store.Add(parcel)
	require.NoError(t, err)
	//и наличии идентификатора
	assert.NotEqual(t, 0, number)

	// set status Registered to Sent
	// обновите статус, убедитесь в отсутствии ошибки
	err = store.SetStatus(number, ParcelStatusSent)
	require.NoError(t, err)
	// check
	// получите добавленную посылку и убедитесь, что статус обновился
	p, err := store.Get(number)
	require.NoError(t, err)
	assert.Equal(t, ParcelStatusSent, p.Status)

	// set status Registered to Delivered
	// обновите статус, убедитесь в отсутствии ошибки
	err = store.SetStatus(number, ParcelStatusDelivered)
	require.NoError(t, err)
	// check
	p, err = store.Get(number)
	require.NoError(t, err)
	assert.Equal(t, ParcelStatusDelivered, p.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	store := NewParcelStore(db)
	defer db.Close()

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
		// добавьте новую посылку в БД, убедитесь в отсутствии ошибки
		number, err := store.Add(parcels[i])
		require.NoError(t, err)
		//и наличии идентификатора
		assert.NotEqual(t, 0, number)

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = number
		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[number] = parcels[i]
	}

	// get by client
	// получите список посылок по идентификатору клиента, сохранённого в переменной client
	storedParcels, err := store.GetByClient(client)
	// убедитесь в отсутствии ошибки
	require.NoError(t, err)
	// убедитесь, что количество полученных посылок совпадает с количеством добавленных
	assert.Equal(t, 3, len(storedParcels))

	// check
	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		expected, ok := parcelMap[parcel.Number]
		require.True(t, ok)

		// убедитесь, что значения полей полученных посылок заполнены верно
		assert.Equal(t, expected.Number, parcel.Number)
		assert.Equal(t, expected.Client, parcel.Client)
		assert.Equal(t, expected.Status, parcel.Status)
		assert.Equal(t, expected.Address, parcel.Address)
		assert.Equal(t, expected.CreatedAt, parcel.CreatedAt)
	}
}
