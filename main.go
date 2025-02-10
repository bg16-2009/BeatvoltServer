package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/bg16_2009/beatvolt_server/models"

	"github.com/go-chi/jwtauth/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var tokenAuth *jwtauth.JWTAuth

type Code struct {
	code         string
	creationTime time.Time
}

var (
	codeStore  = make(map[string]Code)
	storeMutex sync.Mutex
	db         *gorm.DB
)

func main() {
	loadTemplates()

	var err error

	tokenAuth = jwtauth.New("HS256", []byte(secret_key), nil)

	db, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Error migrating user table: ", err)
	}
	err = db.AutoMigrate(&models.Robot{})
	if err != nil {
		log.Fatal("Error migrating robot table: ", err)
	}

	r := makeRouter()

	log.Println("Server starting on :443")
	err = http.ListenAndServeTLS(":443", certFile, keyFile, r)
	if err != nil {
		log.Fatal(err)
	}
}
