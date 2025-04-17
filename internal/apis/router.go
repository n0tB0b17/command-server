package apis

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type APIServer struct {
	Port int
	s    *http.Server
}

func NewAPIServer() *APIServer {
	return &APIServer{
		Port: 8888,
	}
}

func (a *APIServer) Start() {
	addr := fmt.Sprintf(":%d", a.Port)

	router := mux.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"GET", "POST", "OPTIONS"},
		AllowedMethods: []string{"Content-Type", "Accept", "Content-Length"},
	})

	handler := c.Handler(router)
	a.s = &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	a.s.ListenAndServe()
}
