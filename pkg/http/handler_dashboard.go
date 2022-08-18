package http

import (
	"net/http"
	"text/template"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	gui "github.com/momosbasement/radiomomo/pkg/gui"
	model "github.com/momosbasement/radiomomo/pkg/models"
	"gorm.io/gorm"
)

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/static/pic/logo-momo-black.png")
}

func dashboardIndex(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tracks []model.Track
		db.Find(&tracks)

		p := gui.Page{"Titre de ma page", "track", tracks}
		t := template.New("Home")

		t = template.Must(t.ParseFiles("web/tpl/dashboard/layout.html", "web/tpl/dashboard/menu.html", "web/tpl/dashboard/navbar.html", "web/tpl/dashboard/tracks.html"))
		err := t.ExecuteTemplate(w, "layout", p)

		if err != nil {
			log.Fatalf("Template execution: %s", err)
		}
		log.Printf("/GET /dashboard/index\n")
	})
}

func dashboardUpload(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tracks []model.Track
		db.Find(&tracks)

		p := gui.Page{"Titre de ma page", "track", tracks}
		t := template.New("Home")

		t = template.Must(t.ParseFiles("web/tpl/dashboard/layout.html", "web/tpl/dashboard/menu.html", "web/tpl/dashboard/navbar.html", "web/tpl/dashboard/upload.html"))
		err := t.ExecuteTemplate(w, "layout", p)

		if err != nil {
			log.Fatalf("Template execution: %s", err)
		}
		log.Printf("/GET /dashboard/upload\n")
	})
}

func dashboardTrackView(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var track model.Track

		vars := mux.Vars(r)
		key := vars["id"]
		db.First(&track, key)

		t := template.New("Home")

		t = template.Must(t.ParseFiles("web/tpl/dashboard/layout.html", "web/tpl/dashboard/menu.html", "web/tpl/dashboard/navbar.html", "web/tpl/dashboard/trackview.html"))
		err := t.ExecuteTemplate(w, "layout", track)

		if err != nil {
			log.Fatalf("Template execution: %s", err)
		}
		log.Info("/GET /dashboard/track/view/", key)
	})
}

func dashboardTrackDelete(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var track model.Track
		vars := mux.Vars(r)
		key := vars["id"]
		db.First(&track, key)
		track.Delete()
		db.Delete(&track, key)
		http.Redirect(w, r, "/dashboard/tracks", http.StatusFound)
		log.Printf("/GET /dashboard/track/delete/%s\n", key)
	})
}

func dashboardSettings(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tracks []model.Track
		db.Find(&tracks)

		p := gui.Page{"Titre de ma page", "track", tracks}
		t := template.New("Home")

		t = template.Must(t.ParseFiles("web/tpl/dashboard/layout.html", "web/tpl/dashboard/menu.html", "web/tpl/dashboard/navbar.html", "web/tpl/dashboard/settings.html"))
		err := t.ExecuteTemplate(w, "layout", p)

		if err != nil {
			log.Fatalf("Template execution: %s", err)
		}
		log.Printf("/GET /dashboard/settings\n")
	})
}
