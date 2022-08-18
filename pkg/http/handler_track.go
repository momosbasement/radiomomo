package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	model "github.com/momosbasement/radiomomo/pkg/models"
	"gorm.io/gorm"
)

func createNewTrack(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := ioutil.ReadAll(r.Body)
		var track model.Track
		json.Unmarshal(reqBody, &track)
		db.Create(&track)
		json.NewEncoder(w).Encode(track)
		log.Info("/POST /track", track.ID)
	})
}

func updateTrack(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var track model.Track
		reqBody, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(reqBody, &track)
		db.Save(&track)
		json.NewEncoder(w).Encode(track)
		w.Header().Set("Content-Type", "application/json")
		log.Info("/PUT /track/", track.ID)
	})
}

func deleteTrack(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var track model.Track
		vars := mux.Vars(r)
		key := vars["id"]
		db.Delete(&track, key)
		json.NewEncoder(w).Encode(track)
		w.Header().Set("Content-Type", "application/json")
		log.Info("/DELETE /track/", key)
	})
}

func returnAllTracks(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tracks []model.Track
		db.Find(&tracks)
		json.NewEncoder(w).Encode(tracks)
		w.Header().Set("Content-Type", "application/json")
		log.Info("/GET /tracks")
	})
}

func returnSingleTrack(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var track model.Track
		vars := mux.Vars(r)
		key := vars["id"]
		db.First(&track, key)
		json.NewEncoder(w).Encode(track)
		w.Header().Set("Content-Type", "application/json")
		log.Info("/GET /track/", key)
	})
}

func getCover(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var track model.Track
		vars := mux.Vars(r)
		key := vars["id"]
		db.First(&track, key)
		w.Header().Set("Content-Type", "image/jpeg")
		if track.Cover == "" {
			http.ServeFile(w, r, "./web/static/img/track.png")
		} else {
			http.ServeFile(w, r, track.Cover)
		}
		log.Info("/GET /track/cover/", key)
	})
}

func uploadFile(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("/POST /api/track/upload\n")
		// Parse our multipart form, 10 << 20 specifies a maximum
		// upload of 10 MB files.
		r.ParseMultipartForm(32 << 20)
		// FormFile returns the first file for the given key `myFile`
		// it also returns the FileHeader so we can get the Filename,
		// the Header and the size of the file
		file, handler, err := r.FormFile("track-file")
		if err != nil {
			log.Printf("Error while uploading file")
			return
		}
		defer file.Close()
		log.Info("Uploaded File: %+v File Size: %+v MIME Header: %+v", handler.Filename, handler.Size, handler.Header)

		tempFile, err := ioutil.TempFile(os.Getenv("RADIO_MUSIC_DIR"), "*.mp3")
		if err != nil {
			fmt.Println(err)
		}
		defer tempFile.Close()

		// read all of the contents of our uploaded file into a
		// byte array
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}
		// write this byte array to our temporary file
		tempFile.Write(fileBytes)
		// return that we have successfully uploaded our file!
		log.Info("Successfully Uploaded File %s\n", file)
		//io.Copy(tempFile, file)
		// save it into database
		var track model.Track
		track.MP3File = tempFile.Name()
		track.SetFieldsWithID3Tags()
		db.Create(&track)
		track.SetCover()
		db.Save(track)
		//json.NewEncoder(w).Encode(track)
		log.Info("/POST /api/upload\n")
		http.Redirect(w, r, "/dashboard/index", http.StatusFound)
	})
}
