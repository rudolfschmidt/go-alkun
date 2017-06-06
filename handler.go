package alkun

import (
	"net/http"
)

type Handler struct {
	processes []Process
}

func (c *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	for _, p := range c.processes {

		next := p.execute(w, r)

		if !next {
			break
		}

	}

}
