package time_tracker

import (
	"log"
	"net/http"
)

type Server struct {
	http.Server
}

func (s *Server) Start() error {
	log.Println("start serving in ", s.Addr)
	return s.ListenAndServe()
}
