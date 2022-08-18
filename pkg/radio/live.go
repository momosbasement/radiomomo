package radio

import (
	"io"
	"net/http"

	httpRouter "github.com/julienschmidt/httprouter"
)

const sampleRate = 44100
const seconds = 1

func handlerLive(w http.ResponseWriter, req *http.Request, _ httpRouter.Params) {

	buf := make([]byte, 4*1024)

	for {
		n, err := req.Body.Read(buf)
		if n > 0 {
			w.Write(buf[:n])
		}

		if err != nil {
			if err == io.EOF {
				w.Header().Set("Status", "200 OK")
				req.Body.Close()
			}
			break
		}
	}
}
