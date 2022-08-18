package models

import (
	"errors"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	//"github.com/tcolgate/mp3"

	"github.com/bogem/id3v2"

	mp3 "github.com/hajimehoshi/go-mp3"
)

type Track struct {
	Model
	Artist  string `gorm:"type:varchar(50)" json:"artist" validate:"required"`
	Name    string `gorm:"type:varchar(50)" json:"name" validate:"required"`
	Length  int64  `gorm:"type:int(4)" json:"length" validate:"required"`
	Label   string `gorm:"type:varchar(50)" json:"label"`
	Genre   string `gorm:"type:varchar(50)" json:"genre"`
	Album   string `gorm:"type:varchar(100)" json:"album"`
	Year    string `gorm:"type:int(4)" json:"year" validate:"required,year"`
	MP3File string `gorm:"type:varchar(256)" json:"mp3file"`
	Cover   string `gorm:"type:varchar(256)" json:"cover"`
}

type Tracks struct {
	Playlist []Track
}

func (t *Track) SetFieldsWithID3Tags() *Track {
	tag, err := id3v2.Open(t.MP3File, id3v2.Options{Parse: true})
	if err != nil {
		log.Fatal("Error while opening mp3 file: ", err)
	}
	defer tag.Close()
	var result = ""
	if result = tag.Artist(); result == "" {
		result = "Unknown"
	}
	t.Artist = result
	if result = tag.Title(); result == "" {
		result = "Unknown"
	}
	t.Name = result
	if result = tag.Genre(); result == "" {
		result = "Unknown"
	}
	t.Genre = result
	if result = tag.Album(); result == "" {
		result = "Unknown"
	}
	t.Album = result
	if result = tag.Genre(); result == "" {
		result = "Unknown"
	}
	if result = tag.Year(); result == "" {
		result = "Unknown"
	}
	t.Year = result

	t.Length, _ = t.GetLength()
	return t
}

func (t *Track) TableName() string {
	return "track"
}

func (t *Track) Delete() {
	if _, err := os.Stat(t.MP3File); errors.Is(err, os.ErrNotExist) {
		log.Warn("The MP3 file %s does not exists", t.MP3File)
	} else {
		e := os.Remove(t.MP3File)
		if e != nil {
			log.Warn(e)
		}
	}
	if _, err := os.Stat(t.Cover); errors.Is(err, os.ErrNotExist) {
		log.Warn("The MP3 file " + t.MP3File + " does not exists")
	} else {
		e := os.Remove(t.Cover)
		if e != nil {
			log.Warn("The MP3 cover file " + t.Cover + " does not exists")
		}
	}
	log.Info("Track files %d sucessfully deleted", t.ID)
}

func (t *Track) SetCover() {
	tag, err := id3v2.Open(t.MP3File, id3v2.Options{Parse: true})
	if err != nil {
		log.Fatal("Error while opening mp3 file: ", err)
	}
	defer tag.Close()
	// Manage Picture
	pictures := tag.GetFrames(tag.CommonID("Attached picture"))
	for _, f := range pictures {
		pic, ok := f.(id3v2.PictureFrame)
		if !ok {
			log.Warn("Couldn't assert picture frame")
		}
		//TODO path in paramater
		var coverPath = ""
		if len(pic.Picture) > 0 {
			coverPath = "./web/static/img/cover/" + strconv.FormatUint(uint64(t.ID), 10) + ".jpeg"
			err = os.WriteFile(coverPath, pic.Picture, 0644)
			log.Info("Write cover from ID3Tag")
		} else {
			coverPath = "./web/static/img/track.png"
			log.Warn("No cover")
		}

		t.Cover = coverPath
	}
}

func (t *Track) GetLength() (int64, error) {
	f, err := os.Open(t.MP3File)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		return 0, err
	}
	log.Printf("Track length %d", d.Length())
	return d.Length(), nil
}
