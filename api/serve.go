package api

import (
    "net/http"
)

type Route func(http.ResponseWriter, *http.Request)

type Server struct {
	mux *http.ServeMux
}

func NewServer() *Server{
	return &Server{
		mux: http.NewServeMux(),
	}
}

func (s *Server) Register(uri string, fn Route){
    s.mux.HandleFunc(uri, fn);
}

func (s *Server) RegisterAllRoutes(){
  s.Register("/create",Create());
}