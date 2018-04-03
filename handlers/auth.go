package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gopicluster/models"
)

func (env *Env) AppMiddleware(h models.AppHandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if e := h(w, r); e != nil {
			log.Error(e.Message + "-" + e.Error.Error())
			http.Error(w, e.Message, e.Code)
		}
	})
}

func (env *Env) LogingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Debug(req.RequestURI)
		h.ServeHTTP(w, req)
	})
}

func (env *Env) HeadersMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, req)
	})
}
