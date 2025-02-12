package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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

func showHelp(name string) {
	fmt.Printf(string(`Subcommands:
- %s promote <username> : makes <username> an admin
- %s add_robot : add a new robot
- %s start : Starts the server(this is the default behaviour of running with no subcommand)
`), name, name, name)
}

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

	args := os.Args

	if len(args) > 1 && args[1] != "start" {
		if args[1] == "help" {
			showHelp(args[0])
		} else if args[1] == "promote" {
			if len(args) == 1 {
				fmt.Printf("No username provided\n")
				return
			}
			var u models.User
			err := db.Where("username = ?", args[2]).First(&u).Error
			if err != nil {
				fmt.Printf("Error finding user: %v\n", err)
				return
			}
			err = db.Model(&u).Update("is_admin", true).Error
			if err != nil {
				fmt.Printf("Error making admin: %v\n", err)
			} else {
				fmt.Printf("User promoted successfully!\n")
			}
			return
		} else if args[1] == "add_robot" {
			var username, password string
			fmt.Printf("Enter username for new robot: ")
			fmt.Scanf("%s", &username)
			fmt.Printf("Enter password for new robot: ")
			fmt.Scanf("%s", &password)

			bot := models.Robot{Username: username}
			err := bot.SetPassword(password)
			if err != nil {
				fmt.Printf("Error encrypting password: %v\n", err)
				return
			}
			err = db.Create(&bot).Error
			if err != nil {
				fmt.Printf("Error creating robot: %v\n", err)
				return
			}
			fmt.Printf("Robot %s created successfully\n", username)
		} else {
			fmt.Printf("Invalid subcomamnd %s. Use \"%s help\" to show the help.\n", args[1], args[0])
		}
	} else {
		r := makeRouter()

		log.Println("Server starting on :443")
		err = http.ListenAndServeTLS(":443", certFile, keyFile, r)
		if err != nil {
			log.Fatal(err)
		}
	}
}
