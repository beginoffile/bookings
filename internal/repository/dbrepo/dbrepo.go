package dbrepo

import (
	"database/sql"

	"github.com/beginoffile/bookings/internal/config"
	"github.com/beginoffile/bookings/internal/repository"
)

type postgressDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgressRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgressDBRepo{
		App: a,
		DB:  conn,
	}
}

func NewTestingRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}
