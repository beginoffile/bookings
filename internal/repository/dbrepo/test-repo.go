package dbrepo

import (
	"errors"
	"time"

	"github.com/beginoffile/bookings/internal/models"
)

func (m *testDBRepo) AllUser() bool {
	return true
}

// InsertReservation insert a reservation into a database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {

	if res.RoomID == 2 {
		return 0, errors.New("Some error")
	}

	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *testDBRepo) InsertRoomRestriction(res models.RoomRestriction) error {

	if res.RoomID == 1000 {
		return errors.New("some Error")
	}
	return nil

}

// SearchAvailabilityByDatesByRoomID returns true if availability exists form RoomID, and false if no avaliability exists
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, RoomID int) (bool, error) {

	return false, nil

}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any, for given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil

}

// GetRoomByID gets a room by id
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {

	var room models.Room

	if id > 2 {
		return room, errors.New("Some error")

	}

	return room, nil

}

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var user models.User
	return user, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 1, "", nil
}
