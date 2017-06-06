package alkun

import (
	"strings"
	"net/http"
)

const (
	SLASH = "/"
	PARAM_DELIMITER = ":"
)

type Process interface {
	execute(w http.ResponseWriter, r *http.Request) bool
}

type UseProcess struct {
	route route
}

type PathProcess struct {
	path  string
	route route
}

type MethodPathProcess struct {
	method string
	path   string
	route  route
}

func (c *UseProcess) execute(w http.ResponseWriter, r *http.Request) bool {

	req := new(Request)
	req.HttpRequest = r

	res := new(Response)
	res.Writer = w

	c.route(req, res)

	return res.next

}

func (c *PathProcess) execute(w http.ResponseWriter, r *http.Request) bool {

	if !acceptedPath(r.URL.Path, c.path) {
		return true
	}

	req := new(Request)
	req.HttpRequest = r
	req.Path = c.path

	res := new(Response)
	res.Writer = w

	c.route(req, res)

	return res.next

}

func (c *MethodPathProcess) execute(w http.ResponseWriter, r *http.Request) bool {

	if !acceptMethod(r.Method, c.method) {
		return true
	}

	if !acceptedPath(r.URL.Path, c.path) {
		return true
	}

	req := new(Request)
	req.HttpRequest = r
	req.Path = c.path

	res := new(Response)
	res.Writer = w

	c.route(req, res)

	return res.next

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