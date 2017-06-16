package alkun

import (
	"log"
	"net/http"
)

type Server struct {
	processes []Process
}

func (s *Server) Use(route route) {
	s.processes = append(s.processes, &PlainProcess{route})
}

func (s *Server) Filter(path string, route route) {
	s.processes = append(s.processes, &PathProcess{path, route})
}

func (s *Server) Get(path string, route route) {
	s.processes = append(s.processes, &MethodPathProcess{http.MethodGet, path, route})
}

func (s *Server) Post(path string, route route) {
	s.processes = append(s.processes, &MethodPathProcess{http.MethodPost, path, route})
}

func (s *Server) Put(path string, route route) {
	s.processes = append(s.processes, &MethodPathProcess{http.MethodPut, path, route})
}

func (s *Server) Delete(path string, route route) {
	s.processes = append(s.processes, &MethodPathProcess{http.MethodDelete, path, route})
}

func (s *Server) Head(path string, route route) {
	s.processes = append(s.processes, &MethodPathProcess{http.MethodHead, path, route})
}

func (s *Server) Options(path string, route route) {
	s.processes = append(s.processes, &MethodPathProcess{http.MethodOptions, path, route})
}

func (s *Server) Start(port string) {

	err := http.ListenAndServe(port, &Handler{processes: s.processes})

	if err != nil {
		log.Fatalf("error: %s", err)
	}

}