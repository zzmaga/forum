package web

import (
	"fmt"
	"log"
	"net/http"
)

// TODO: add from env.configs
type Configs struct {
	Port string
	// ReadTimeout  int
	// WriteTimeOut int
	// IdleTimeout int
	// MaxHeaderBytes int
}

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(configs *Configs, handler http.Handler) error {
	configs.Port = "8080"
	s.httpServer = &http.Server{
		Addr:    configs.Port,
		Handler: handler,
		// MaxHeaderBytes: configs.MaxHeaderBytes,
		// ReadTimeout:    time.Duration(configs.ReadTimeout * int(time.Millisecond)),
		// WriteTimeout:   time.Duration(configs.WriteTimeout * int(time.Millisecond)),
		// IdleTimeout:    time.Duration(configs.IdleTimeout * int(time.Millisecond)),
	}

	log.Printf("Server runs on http://localhost%s\n", s.httpServer.Addr)
	err := s.httpServer.ListenAndServe()
	return fmt.Errorf("Run: %w", err)
}
