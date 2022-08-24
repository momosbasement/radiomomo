package main

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	database "github.com/momosbasement/radiomomo/pkg/db"
	http "github.com/momosbasement/radiomomo/pkg/http"
	radio "github.com/momosbasement/radiomomo/pkg/radio"
	log "github.com/sirupsen/logrus"
)

// ------ MAIN

func main() {
	// LOAD ENV VARS
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// RADIO
	radio := new(radio.Radio)
	port, _ := strconv.Atoi(os.Getenv("RADIO_PORT"))
	radio.Init(port, os.Getenv("RADIO_MUSIC_DIR"))
	go radio.PlayRadio()
	//DB
	conn := new(database.Connection)
	conn.Init(os.Getenv("RADIO_DB_DSN"))
	radio.SetDatabase(conn.GetConnection())
	radio.GeneratePlaylist()
	//HANDLER
	http.HandleRequests(radio)
}
