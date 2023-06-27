package repository

import (
	"time"

	"github.com/beginoffile/bookings/cmd/internal/models"
)

type DatabaseRepo interface {
	AllUser() bool

	InsertReservation(res models.Reservation) (int, error)

	InsertRoomRestriction(res models.RoomRestriction) error

	SearchAvailabilityByDatesByRoomID(start, end time.Time, RoomID int) (bool, error)

	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
}
