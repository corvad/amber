package api

import (
	"net/http"
)

func Create() Route {
  return func(w http.ResponseWriter, r *http.Request){
    if(r.Method != "POST"){
      w.
	}
  }
}