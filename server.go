package alkun

import (
	"log"
	"net/http"
)

type Server struct {
	processes []process
}

func (s *Server) Use(route default_route) {
	s.processes = append(s.processes, &plain_process{route})
}

func (s *Server) Filter(path string, route default_route) {
	s.processes = append(s.processes, &path_process{path, route})
}

func (s *Server) Get(path string, route default_route) {
	s.processes = append(s.processes, &method_path_process{http.MethodGet, path, route})
}

func (s *Server) Post(path string, route default_route) {
	s.processes = append(s.processes, &method_path_process{http.MethodPost, path, route})
}

func (s *Server) Put(path string, route default_route) {
	s.processes = append(s.processes, &method_path_process{http.MethodPut, path, route})
}

func (s *Server) Delete(path string, route default_route) {
	s.processes = append(s.processes, &method_path_process{http.MethodDelete, path, route})
}

func (s *Server) Head(path string, route default_route) {
	s.processes = append(s.processes, &method_path_process{http.MethodHead, path, route})
}

func (s *Server) Options(path string, route default_route) {
	s.processes = append(s.processes, &method_path_process{http.MethodOptions, path, route})
}

func (s *Server) Exception(route exception_route) {
	s.processes = append(s.processes, &error_process{route})
}

func (s *Server) Start(port string) {

	err := http.ListenAndServe(port, handle_processes(s.processes))

	if err != nil {
		log.Fatalf("error: %s", err)
	}

}

func handle_processes(processes []process) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		handle_processes_indexed(processes, 0, w, r)

	}

}

func handle_processes_indexed(processes []process, index int, w http.ResponseWriter, r *http.Request) {

	//handle runtime exceptions
	defer func() {
		if re := recover(); re != nil {
			handle_error_routes(processes, index+1, re, w, r)
		}
	}()

	if index >= len(processes) {
		return
	}

	switch p := processes[index].(type) {

	case default_process:

		next, err := p.execute(w, r)

		if err != nil {
			handle_error_routes(processes, index+1, err, w, r)
			return
		}

		if !next {
			return
		}

	}

	handle_processes_indexed(processes, index+1, w, r)

}

func handle_error_routes(processes []process, index int, err interface{}, w http.ResponseWriter, r *http.Request) {

	//handle runtime exceptions (errors in error handling routes)
	defer func() {
		if re := recover(); re != nil {
			handle_error_routes(processes, index+1, re, w, r)
		}
	}()

	if index >= len(processes) {
		return
	}

	switch p := processes[index].(type) {

	case *error_process:

		next, new_err := p.execute(err, w, r)

		if new_err != nil {
			handle_error_routes(processes, index+1, new_err, w, r)
			return
		}

		if !next {
			return
		}

	}

	handle_error_routes(processes, index+1, err, w, r)

}
