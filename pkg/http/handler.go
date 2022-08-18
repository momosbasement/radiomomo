package http

import (
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	radio "github.com/momosbasement/radiomomo/pkg/radio"
)

func HandleRequests(radio *radio.Radio) {
	myRouter := mux.NewRouter().StrictSlash(true)

	s := http.StripPrefix("/web/static/", http.FileServer(http.Dir("web/static")))
	// DASHBOARD
	myRouter.PathPrefix("/web/static/").Handler(s)
	myRouter.Handle("/", LivePage(radio, radio.GetDatabase()))
	myRouter.Handle("/dashboard/track/view/{id}", dashboardTrackView(radio.GetDatabase()))
	myRouter.Handle("/dashboard/track/delete/{id}", dashboardTrackDelete(radio.GetDatabase()))
	myRouter.Handle("/dashboard/index", dashboardIndex(radio.GetDatabase()))
	myRouter.Handle("/dashboard/tracks", dashboardIndex(radio.GetDatabase()))
	myRouter.Handle("/dashboard/upload", dashboardUpload(radio.GetDatabase()))
	myRouter.Handle("/dashboard/settings", dashboardSettings(radio.GetDatabase()))
	myRouter.HandleFunc("/favicon.ico", faviconHandler)
	// RADIO RESOURCE
	myRouter.Handle("/api/radio/next", nextTrack(radio)).Methods("GET")
	// TRACK API
	myRouter.Handle("/api/track/play/{id}", radio.StreamTrack(radio.GetDatabase())).Methods("GET")
	myRouter.Handle("/api/track/cover/{id}", getCover(radio.GetDatabase())).Methods("GET")
	myRouter.Handle("/api/track/upload", uploadFile(radio.GetDatabase())).Methods("POST")
	myRouter.Handle("/api/track/{id}", deleteTrack(radio.GetDatabase())).Methods("DELETE")
	myRouter.Handle("/api/track/{id}", returnSingleTrack(radio.GetDatabase())).Methods("GET")
	myRouter.Handle("/api/tracks", returnAllTracks(radio.GetDatabase()))
	myRouter.Handle("/api/track", updateTrack(radio.GetDatabase())).Methods("PUT")
	myRouter.Handle("/api/track", createNewTrack(radio.GetDatabase())).Methods("POST")
	// STREAM
	myRouter.Handle("/live", handleLive(radio))

	log.Info("Serving on HTTP port: ", radio.Port)
	_port := strconv.Itoa(radio.Port)
	log.Info(http.ListenAndServe(":"+_port, myRouter))
}
