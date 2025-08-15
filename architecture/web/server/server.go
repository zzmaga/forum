package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Configs struct {
	Port           string `cenv:"port"`
	ReadTimeout    int    `cenv:"read_timeout_ms"`
	WriteTimeout   int    `cenv:"write_timeout_ms"`
	IdleTimeout    int    `cenv:"idle_timeout_ms"`
	MaxHeaderBytes int    `cenv:"max_header_bytes"`
}

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(configs *Configs, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           configs.Port,
		Handler:        handler,
		MaxHeaderBytes: configs.MaxHeaderBytes,
		ReadTimeout:    time.Duration(configs.ReadTimeout * int(time.Millisecond)),
		WriteTimeout:   time.Duration(configs.WriteTimeout * int(time.Millisecond)),
		IdleTimeout:    time.Duration(configs.IdleTimeout * int(time.Millisecond)),
	}

	log.Printf("Server runs on http://localhost%s\n", s.httpServer.Addr)
	err := s.httpServer.ListenAndServe()
	return fmt.Errorf("Run: %w", err)
}
