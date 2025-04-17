package apis

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/n0tB0b17/isner/internal/connections"
	"github.com/n0tB0b17/isner/internal/storage"
	"github.com/rs/cors"
)

type APIServer struct {
	Port      int
	s         *http.Server
	wsHandler *connections.WebSocketHandler
	store     storage.Repository
}

func NewAPIServer(ws *connections.WebSocketHandler, store storage.Repository) *APIServer {
	return &APIServer{
		Port:      8888,
		wsHandler: ws,
		store:     store,
	}
}

func (a *APIServer) Start() error {
	addr := fmt.Sprintf(":%d", a.Port)

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/connect", a.wsHandler.ServeHttp)
	// add other endpoint
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

	fmt.Printf("Starting API server on address: %s \n", addr)
	return a.s.ListenAndServe()
}

func (a *APIServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.s.Shutdown(ctx); err != nil {
		return fmt.Errorf("error while shutting down API server: %v", err)
	}

	return nil
}
