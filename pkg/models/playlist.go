package models

import (
	"math/rand"
	"time"
)

type Playlist struct {
	Tracklist []Track
}

func (p *Playlist) GenerateRandomTracks() {

}

func (p *Playlist) GetRandomTrack() Track {
	rand.Seed(time.Now().UnixNano())
	return p.Tracklist[rand.Intn(len(p.Tracklist))]
}
