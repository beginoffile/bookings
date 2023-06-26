package dbrepo

import (
	"context"
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
