package webserver

import (
	"log"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type WebServer struct {
	Router        http.ServeMux
	Hanlders      map[string]http.Handler
	WebServerPort string
}

func New(port string) *WebServer {
	return &WebServer{
		Hanlders:      make(map[string]http.Handler),
		WebServerPort: port,
	}
}

func (w *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	w.Hanlders[path] = otelhttp.NewHandler(http.HandlerFunc(handler), path)
}

func (w *WebServer) Start(port string) {
	for path, handler := range w.Hanlders {
		w.Router.Handle(path, handler)
	}

	log.Println("Starting web server...")
	if err := http.ListenAndServe(port, &w.Router); err != nil {
		panic(err)
	}
}
