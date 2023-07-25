package dbrepo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/beginoffile/bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
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

// GetUserByID returns a user by id
func (m *postgressDBRepo) GetUserByID(id int) (models.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `Select id, first_name, last_name, email, password, access_level, created_at, update_at
	from public.user
	Where id = $1`

	var u models.User

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.AccessLevel, &u.CreateAt, &u.UpdateAt)

	if err != nil {
		return u, err
	}

	return u, nil

}

// UpdateUser updates a user in the database
func (m *postgressDBRepo) UpdateUser(u models.User) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `Update public.user 
	Set first_name = $1, 
		last_name = $2, 
		email= $3, 		
		access_level =$4, 		
		update_at=$5	
	Where id = $6`

	_, err := m.DB.ExecContext(ctx, query, u.FirstName, u.LastName, u.Email, u.AccessLevel, time.Now())

	if err != nil {
		return err
	}

	return nil

}

// Authenticate authenticate a user
func (m *postgressDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var id int
	var hashedPassword string

	query := `Select id, password from public.user Where email = $1`

	err := m.DB.QueryRowContext(ctx, query, email).Scan(&id, &hashedPassword)

	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("Incorrect Password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil

}

// AllReservations returns a slice of all reservations
func (m *postgressDBRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var reservations []models.Reservation

	query := `select t1.id, t1.first_name, t1.last_name, t1.email, t1.phone, t1.start_date, t1.end_date, t1.room_id, t1.created_at, t1.updated_at, t1.processed, t2.id, t2.room_name
	from reservation t1
		Left Join room t2
	       On t2.id = t1.room_id
	Order by t1.start_date
		`

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return reservations, err
	}

	defer rows.Close()

	for rows.Next() {
		var i models.Reservation

		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreateAt,
			&i.UpdateAt,
			&i.Processed,
			&i.Room.ID,
			&i.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}

		reservations = append(reservations, i)

	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil

}

// AllNewReservationWhere t1.processed = 0 returns a slice of all reservations
func (m *postgressDBRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var reservations []models.Reservation

	query := `select t1.id, t1.first_name, t1.last_name, t1.email, t1.phone, t1.start_date, t1.end_date, t1.room_id, t1.created_at, t1.updated_at, t2.id, t2.room_name
	from reservation t1
		Left Join room t2
	       On t2.id = t1.room_id
	Where t1.processed = 0
	Order by t1.start_date
		`

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return reservations, err
	}

	defer rows.Close()

	for rows.Next() {
		var i models.Reservation

		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreateAt,
			&i.UpdateAt,
			&i.Room.ID,
			&i.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}

		reservations = append(reservations, i)

	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil

}

// GetReservationByID returns  one reservation by ID
func (m *postgressDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var res models.Reservation

	query := `select t1.id, t1.first_name, t1.last_name, t1.email, t1.phone, t1.start_date, t1.end_date, t1.room_id, t1.created_at, t1.updated_at, t1.processed, t2.id, t2.room_name
	from reservation t1
		Left Join room t2
	       On t2.id = t1.room_id
	Where t1.id = $1
		`

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Phone,
		&res.StartDate,
		&res.EndDate,
		&res.RoomID,
		&res.CreateAt,
		&res.UpdateAt,
		&res.Processed,
		&res.Room.ID,
		&res.Room.RoomName,
	)

	if err != nil {
		return res, err
	}

	return res, nil

}

// UpdateReservation updates a reservation in the database
func (m *postgressDBRepo) UpdateReservation(u models.Reservation) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `Update reservation
	Set first_name = $1, 
		last_name = $2, 
		email= $3, 		
		phone=$4,
		updated_at=$5
	Where id = $6`

	_, err := m.DB.ExecContext(ctx, query, u.FirstName, u.LastName, u.Email, u.Phone, time.Now(), u.ID)

	if err != nil {
		return err
	}

	return nil

}

// DeleteReservation deletes one reservation by id
func (m *postgressDBRepo) DeleteReservation(id int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `Delete from reservation	
	Where id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil

}

// UpdateProcessedForReservation updates processed for a reservation by id
func (m *postgressDBRepo) UpdateProcessedForReservation(id, processed int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `Update reservation
	Set processed = $1
	Where id = $2`

	_, err := m.DB.ExecContext(ctx, query, processed, id)

	if err != nil {
		return err
	}

	return nil

}

// AllRooms Get all rooms
func (m *postgressDBRepo) AllRooms() ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var rooms []models.Room

	query := `select t1.id, t1.room_name, t1.created_at, t1.updated_at
	from room t1	
	order by t1.room_name
		`

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return rooms, err
	}
	defer rows.Close()

	for rows.Next() {
		var rm models.Room
		err := rows.Scan(
			&rm.ID,
			&rm.RoomName,
			&rm.CreateAt,
			&rm.UpdateAt,
		)

		if err != nil {
			return rooms, err
		}

		rooms = append(rooms, rm)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil

}

// GetRestrictionsForRoomByDate	returns restrictions for a room by date range
func (m *postgressDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var restrictions []models.RoomRestriction

	query := `select t1.id, coalesce(t1.reservation_id,0), t1.restriction_id, t1.room_id, t1.start_date, t1.end_date
	from room_restriction t1	
	Where t1.end_date > $1
	 And t1.start_date <= $2
	 And t1.room_id = $3	
		`

	rows, err := m.DB.QueryContext(ctx, query, start, end, roomID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rr models.RoomRestriction
		err := rows.Scan(
			&rr.ID,
			&rr.ReservationID,
			&rr.RestrictionID,
			&rr.RoomID,
			&rr.StartDate,
			&rr.EndDate,
		)

		if err != nil {
			return nil, err
		}

		restrictions = append(restrictions, rr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return restrictions, nil

}

// InsertBlockForRoom inserts a room restriction
func (m *postgressDBRepo) InsertBlockForRoom(id int, startDate time.Time) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `Insert into room_restriction
	(start_date, end_date, room_id, restriction_id, created_at, updated_at)
	values
	($1, $2, $3, $4, $5, $6)
	`

	_, err := m.DB.ExecContext(ctx, query, startDate, startDate.AddDate(0, 0, 1), id, 2, time.Now(), time.Now())

	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

// DeleteBlockForRoom deletes a room restriction
func (m *postgressDBRepo) DeleteBlockByID(id int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `Delete from room_restriction
	Where id = $1
	`

	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}
