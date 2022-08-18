package http

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"text/template"

	log "github.com/sirupsen/logrus"

	gui "github.com/momosbasement/radiomomo/pkg/gui"
	model "github.com/momosbasement/radiomomo/pkg/models"
	radio "github.com/momosbasement/radiomomo/pkg/radio"
	"gorm.io/gorm"
)

func nextTrack(radio *radio.Radio) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		radio.Next()
		json.NewEncoder(w).Encode(radio.GetTrackOnAir())
		w.Header().Set("Content-Type", "application/json")
		log.Info("/PUT /track\n")
	})
}

func handleLive(radio *radio.Radio) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "audio/mpeg")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Server", "radio-momo")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Transfer-Encoding", "chunked")
		w.Header().Set("X-Content-Type-Options", "nosniff")

		w.WriteHeader(http.StatusOK)

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.NotFound(w, r)
			return
		}
		flusher.Flush()

		//TODO infinite loop
		file, err := os.Open(radio.GetTrackOnAir())
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		p := make([]byte, 44000)
		for {
			n, err := file.Read(p)
			if err == io.EOF {
				break
			}
			w.Write(p[:n])
			flusher.Flush()
		}
		log.Info("/GET /live\n")
	})
}

func LivePage(radio *radio.Radio, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tracks []model.Track
		db.Find(&tracks)

		p := gui.Page{"Radio MOMO", "live", tracks}
		t := template.New("Live Page")
		t = template.Must(t.ParseFiles("web/tpl/radio/layout.html", "web/tpl/radio/header.html", "web/tpl/radio/home.html"))
		err := t.ExecuteTemplate(w, "layout", p)

		if err != nil {
			log.Fatalf("Template execution: %s", err)
		}
		log.Info("/GET /\n")
	})
}
