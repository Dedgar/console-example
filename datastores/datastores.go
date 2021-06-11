package datastores

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/dedgar/console-example/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	// PostMap containes the names of eligible posts and their paths
	PostMap      = make(map[string]string)
	CookieSecret string
	OAuthID      string
	OAuthKey     string
	dbHost       string
	dbPort       string
	dbUser       string
	dbPass       string
	dbName       string
	Subject      string
	CharSet      string
	Sender       string
	Recipient    string
	AuthMap      map[string]bool
	DB           = &gorm.DB{}
	VidMap       = map[string]map[string]map[string]struct{}{
		"divisionrune": map[string]map[string]struct{}{
			"1": map[string]struct{}{
				"1": struct{}{},
				"2": struct{}{},
				"3": struct{}{},
			},
		},
	}
)

// checkDB looks for the Users table in the connected DB (if available)
// and creates the table if it does not already exist.
func checkDB() {
	if !DB.Migrator().HasTable(&models.User{}) {
		fmt.Println("Creating users table")
		DB.Migrator().CreateTable(&models.User{})
	}
}

func init() {
	var appSecrets models.AppSecrets

	filePath := "/secrets/dedgar_secrets.json"
	fileBytes, err := ioutil.ReadFile(filePath)

	if err != nil {
		fmt.Println("Error loading secrets json: ", err)
	}

	err = json.Unmarshal(fileBytes, &appSecrets)
	if err != nil {
		fmt.Println("Error Unmarshaling secrets json: ", err)
	}

	CookieSecret = appSecrets.CookieSecret
	OAuthID = appSecrets.GoogleAuthID
	OAuthKey = appSecrets.GoogleAuthKey
	dbPass = appSecrets.PsqlPassword
	dbUser = appSecrets.PsqlUser
	dbPort = appSecrets.PsqlServicePort
	dbName = appSecrets.PsqlDatabase
	dbHost = appSecrets.PsqlServiceHost
	Subject = appSecrets.Subject
	CharSet = appSecrets.CharSet
	Sender = appSecrets.Sender
	AuthMap = appSecrets.AuthMap
	Recipient = appSecrets.Recipient

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)

	DB, err = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		fmt.Println("Error establishing database connection:", err)
	}

	fmt.Println("printing appsecrets")
	fmt.Printf("%+v", appSecrets)

	checkDB()
}
