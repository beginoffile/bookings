package dbrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/beginoffile/bookings/cmd/internal/models"
)

func (m *postgressDBRepo) AllUser() bool {
	return true
}

// InsertReservation insert a reservation into a database
func (m *postgressDBRepo) InsertReservation(res models.Reservation) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var newID int

	stmt :=
		`insert into reservation
	(first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
	values
	($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`

	fmt.Println("AAAAAAAAAAAAQUI======>", res.FirstName, res.LastName, res.Email, res.Phone, res.StartDate, res.EndDate, res.RoomID, time.Now(), time.Now())

	// _, err := m.DB.ExecContext(ctx, stmt, res.FirstName, res.LastName, res.Email, res.Phone, res.StartDate, res.EndDate, res.RoomID, time.Now(), time.Now())
	err := m.DB.QueryRowContext(ctx, stmt, res.FirstName, res.LastName, res.Email, res.Phone, res.StartDate, res.EndDate, res.RoomID, time.Now(), time.Now()).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *postgressDBRepo) InsertRoomRestriction(res models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	stmt :=
		`insert into room_restriction
	(start_date, end_date, room_id, reservation_id, created_at, updated_at,restriction_id)
	values
	($1,$2,$3,$4,$5,$6,$7)`

	_, err := m.DB.ExecContext(ctx, stmt, res.StartDate, res.EndDate, res.RoomID, res.ReservationID, time.Now(), time.Now(), res.RestrictionID)
	// err := m.DB.QueryRowContext(ctx, stmt, res.FirstName, res.LastName, res.Email, res.Phone, res.StartDate, res.EndDate, res.RoomID, time.Now(), time.Now()).Scan(&newID)

	if err != nil {
		return err
	}

	return nil

}

// SearchAvailabilityByDatesByRoomID returns true if availability exists form RoomID, and false if no avaliability exists
func (m *postgressDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, RoomID int) (bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `select count(id)
			From room_restriction
			Where room_id = $1
			  And end_date > $2
			  And start_date < $3`

	var numRows int

	err := m.DB.QueryRowContext(ctx, query, RoomID, start, end).Scan(&numRows)

	if err != nil {
		return false, nil
	}

	if numRows == 0 {
		return false, nil
	}

	return true, nil

}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any, for given date range
func (m *postgressDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var rooms []models.Room

	query := `
			Select t1.id, t1.room_name
			From room t1
			Where Not exists(Select *
							From room_restriction A
							Where A.room_id = t1.id
							And A.end_date > $1
							And A.start_date < $2
							)`

	rows, err := m.DB.QueryContext(ctx, query, start, end)

	if err != nil {
		return rooms, nil
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}

		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil

}

// GetRoomByID gets a room by id
func (m *postgressDBRepo) GetRoomByID(id int) (models.Room, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var room models.Room

	query := `
	Select t1.id, t1.room_name, t1.created_at, t1.updated_at
	From Room t1
	Where t1.id = $1
	`

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&room.ID,
		&room.RoomName,
		&room.CreateAt,
		&room.UpdateAt,
	)

	if err != nil {
		return room, errors.New("[GetRoomByID.QueryRowContext]" + err.Error())
	}

	return room, nil

}
