package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/beginoffile/bookings/internal/config"
	"github.com/beginoffile/bookings/internal/driver"
	"github.com/beginoffile/bookings/internal/handlers"
	"github.com/beginoffile/bookings/internal/helpers"
	"github.com/beginoffile/bookings/internal/models"
	"github.com/beginoffile/bookings/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func AddValue(x, y int) int {
	return x + y
}

func main() {

	db, err := run()

	if err != nil {
		log.Fatal(err)

	}

	defer db.SQL.Close()

	defer close(app.MailChan)

	fmt.Println("Starting mail Listener...")

	listenForMail()

	// msg := models.MailData{
	// 	To:      "john@do.lat",
	// 	From:    "me@here.com",
	// 	Subject: "Hola",
	// 	Content: "",
	// }

	// app.MailChan <- msg

	// from := "me@here.com"

	// auth := smtp.PlainAuth("", from, "", "localhost")

	// err = smtp.SendMail("localhost:1025", auth, from, []string{"you@here.com"}, []byte("Hello World"))
	// if err != nil {
	// 	log.Println(err)
	// }

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	// fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	// http.ListenAndServe(portNumber, nil)

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)

}

func run() (*driver.DB, error) {

	//What am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	// read flags
	inProduction := flag.Bool("production", true, "Application is in production")
	useCache := flag.Bool("cache", true, "Use template cache")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbName := flag.String("dbname", "", "Database name")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPass := flag.String("dbpass", "", "Database password")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbSSL := flag.String("dbssl", "disable", "Database ssl settings(disable, prefer, require)")

	flag.Parse()

	if *dbName == "" || *dbUser == "" {
		fmt.Println("Mising required flags")
		os.Exit(1)
	}

	mailChan := make(chan models.MailData)

	app.MailChan = mailChan

	app.InProduction = *inProduction
	app.UseCache = *useCache

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	// db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=12345678")
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	db, err := driver.ConnectSQL(connectionString)

	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}

	// defer db.SQL.Close()

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("Cannot create template cache", err)
		return nil, err
	}

	app.TemplateCache = tc
	// app.UseCache = false

	repo := handlers.NewRepo(&app, db)

	handlers.NewHandlers(repo)

	render.NewRenderer(&app)

	helpers.NewHelpers(&app)

	return db, nil

}
