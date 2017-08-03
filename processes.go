package alkun

import (
	"net/http"
	"strings"
)

const (
	SLASH           = "/"
	PARAM_DELIMITER = ":"
)

type process interface {
}

type default_process interface {
	execute(w http.ResponseWriter, r *http.Request) (bool, error)
}

type plain_process struct {
	route default_route
}

type path_process struct {
	path  string
	route default_route
}

type method_path_process struct {
	method string
	path   string
	route  default_route
}

type error_process struct {
	route exception_route
}

func (p *error_process) execute(err interface{}, w http.ResponseWriter, r *http.Request) (bool, error) {

	req := new(Request)
	req.HttpRequest = r

	res := new(Response)
	res.Writer = w

	new_err := p.route(err, req, res)

	if new_err != nil {
		return false, new_err
	}

	return res.next, nil

}

func (p *plain_process) execute(w http.ResponseWriter, r *http.Request) (bool, error) {

	req := new(Request)
	req.HttpRequest = r

	res := new(Response)
	res.Writer = w

	err := p.route(req, res)

	if err != nil {
		return false, err
	}

	return res.next, nil

}

func (p *path_process) execute(w http.ResponseWriter, r *http.Request) (bool, error) {

	if !acceptedPath(r.URL.Path, p.path) {
		return true, nil
	}

	req := new(Request)
	req.HttpRequest = r
	req.Path = p.path

	res := new(Response)
	res.Writer = w

	err := p.route(req, res)

	if err != nil {
		return false, err
	}

	return res.next, nil

}

func (p *method_path_process) execute(w http.ResponseWriter, r *http.Request) (bool, error) {

	if !acceptMethod(r.Method, p.method) {
		return true, nil
	}

	if !acceptedPath(r.URL.Path, p.path) {
		return true, nil
	}

	req := new(Request)
	req.HttpRequest = r
	req.Path = p.path

	res := new(Response)
	res.Writer = w

	err := p.route(req, res)

	if err != nil {
		return false, err
	}

	return res.next, nil

}

func acceptMethod(requestMethod string, providedMethod string) bool {

	return requestMethod == providedMethod

}

func acceptedPath(requestedPath string, providedPath string) bool {

	if requestedPath == providedPath {
		return true
	}

	requested := strings.Split(requestedPath, SLASH)
	provided := strings.Split(providedPath, SLASH)

	if len(requested) != len(provided) {
		return false
	}

	for i := 0; (i < len(requested)) && (i < len(provided)); i++ {

		if !strings.HasPrefix(provided[i], PARAM_DELIMITER) && requested[i] != provided[i] {
			return false
		}

	}

	return true

}
