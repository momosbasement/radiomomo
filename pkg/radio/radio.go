package radio

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	model "github.com/momosbasement/radiomomo/pkg/models"
	file "github.com/momosbasement/radiomomo/pkg/utils"
	"gorm.io/gorm"
)

type Radio struct {
	Port       int
	Folder     string
	Seek       int64
	Uptime     int64
	Buffer     []byte
	TrackOnAir string
	Playlist   []model.Track
	db         *gorm.DB
}

// Constructor
func (r *Radio) Init(port int, folder string) {
	rand.Seed(time.Now().UnixNano())
	r.Port = port
	r.Folder = folder
	r.Seek = 0
	r.TrackOnAir = ""
	r.Playlists, _ = file.FilePathWalkDir(r.Folder)
}

// Set the database connection
func (r *Radio) SetDatabase(d *gorm.DB) {
	r.db = d
}

// Return the database connection
func (r *Radio) GetDatabase() *gorm.DB {
	return r.db
}

// Generate a playlist
func (r *Radio) GeneratePlaylist() {
	var tracks []model.Track
	r.db.Find(&tracks)
	r.Playlists
}

// Return a random track
func (r *Radio) getRandomTrack() string {
	return r.Playlists[rand.Intn(len(r.Playlists))]
}

func (r *Radio) setRandomTrack() {
	r.Seek = 0
	r.TrackOnAir = r.Playlists[rand.Intn(len(r.Playlists))]
	log.Info("Setting new track ", r.TrackOnAir)
}

func (r *Radio) SetTrack(t model.Track) {
	r.Seek = 0
	r.TrackOnAir = t.MP3File
	log.Info("Setting track ", r.TrackOnAir)
}

func (r *Radio) GetTrackOnAir() string {
	return r.TrackOnAir
}

func (r *Radio) getSeekPosition() int64 {
	return r.Seek
}

// it's just an alias
func (r *Radio) Next() {

	r.setRandomTrack()
}

func (radio *Radio) StreamTrack(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var track model.Track
		vars := mux.Vars(r)
		key := vars["id"]

		if result := db.First(&track, key); result.Error != nil {
			http.NotFound(w, r)
			return
		}

		file, err := os.Open(track.MP3File)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "audio/mpeg")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Server", "radio-momo")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Transfer-Encoding", "chunked")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusOK)

		radio.SetTrack(track)

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.NotFound(w, r)
			return
		}
		flusher.Flush()

		//TODO infinite loop

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
		log.Info("/GET /api/track/play/", track.ID)
	})
}

func (r *Radio) PlayRadio() {
	r.setRandomTrack()
	log.Info("Starting playing track " + r.GetTrackOnAir())
	f, err := os.Open(r.GetTrackOnAir())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	//streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	var sampleRate = 44000
	var seconds = 1
	p := make([]byte, sampleRate*seconds)
	for {
		_, err := f.Read(p)
		if err == io.EOF {
			break
		}
		r.Seek++
		r.Uptime++
		//log.Println("Playing .... Seek position ", r.Seek)
		time.Sleep(1 * time.Second)
	}
	// recursivly call to infinite random playing
	log.Info("Finishing playing track " + r.GetTrackOnAir())
	go r.PlayRadio()
}
