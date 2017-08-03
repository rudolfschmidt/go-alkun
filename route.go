package alkun

import (
	"net/http"
	"strings"
)

type default_route func(*Request, *Response) error
type exception_route func(interface{}, *Request, *Response) error
type parameters map[string]string

type Request struct {
	HttpRequest *http.Request
	Path        string
	parameters  parameters
}

func (r *Request) Param(id string) string {

	req := strings.Split(r.HttpRequest.URL.Path, SLASH)
	prov := strings.Split(r.Path, SLASH)

	if r.parameters == nil {

		r.parameters = make(parameters)

		for i := 0; (i < len(req)) && (i < len(prov)); i++ {

			if strings.HasPrefix(prov[i], PARAM_DELIMITER) {
				r.parameters[prov[i]] = req[i]
			}

		}

	}

	if !strings.HasPrefix(id, PARAM_DELIMITER) {
		id = PARAM_DELIMITER + id
	}

	return r.parameters[id]
}

type Response struct {
	Writer http.ResponseWriter
	next   bool
}

func (r *Response) Next() {
	r.next = true
}
