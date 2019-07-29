package service

import "net/http"

func (s *service) SetupHandlers() *http.ServeMux {
	return s.setupHandlers()
}
