package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gopicluster/models"
	"os"
	"os/exec"
)

var (
	Ch   chan models.Resp
	Wsrw *models.WSReaderWriter
	T    time.Time
)

func (env *Env) TestHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	w.Write([]byte(T.String()))
	return nil
}

func (env *Env) CrackHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	type p struct {
		Password string `json:"password"`
		Nbnodes  string `json:"nbnodes"`
	}
	var pwd p
	if err := r.ParseForm(); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "form parsing error",
			Code:    http.StatusBadRequest}
	}
	var decoder = schema.NewDecoder()
	if err := decoder.Decode(&pwd, r.PostForm); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "form decoding error",
			Code:    http.StatusBadRequest}
	}
	if pwd.Nbnodes == "0" {
		pwd.Nbnodes = "1"
	}
	log.WithFields(log.Fields{"pwd.Password": pwd.Password, "pwd.Nbnodes": pwd.Nbnodes}).Debug()

	cmd := exec.Command("/opt/slurm/launch_to_nodes.sh", pwd.Password, pwd.Nbnodes)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Error(err.Error())
		return &models.AppError{
			Error:   err,
			Message: "shell exec command error",
			Code:    http.StatusInternalServerError}
	}
	if err != nil {
		log.Error(err.Error())
		return &models.AppError{
			Error:   err,
			Message: "shell exec command error",
			Code:    http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	return nil
}

// HomeHandler serve the main page
func (env *Env) HomeHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	if e := env.Templates["home"].ExecuteTemplate(w, "base", env); e != nil {
		return &models.AppError{
			Error:   e,
			Code:    http.StatusInternalServerError,
			Message: "error executing template base",
		}
	}

	return nil
}

func (env *Env) StartHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	var (
		jobid string
		node  string
	)

	rquery := r.URL.Query()
	jobid = rquery.Get("jobid")
	node = rquery.Get("node")
	log.WithFields(log.Fields{"jobid": jobid, "node": node}).Debug()

	Ch <- models.Resp{Jobid: jobid, Node: node}

	w.Write([]byte("ok"))
	return nil
}

func (env *Env) StopHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	var (
		jobid  string
		node   string
		result string
		pass   string
	)

	rquery := r.URL.Query()
	jobid = rquery.Get("jobid")
	node = rquery.Get("node")
	result = rquery.Get("result")
	pass = rquery.Get("pass")

	Ch <- models.Resp{Jobid: jobid, Node: node, Result: result, Pass: pass}

	w.Write([]byte("ok"))
	return nil
}

func (env *Env) SocketHandler(w http.ResponseWriter, r *http.Request) {
	// opening the websocket
	Wsrw = models.NewWSsender(w, r)

	for {

		r := Wsrw.Read()

		if r.Nbnodes == "0" {
			r.Nbnodes = "1"
		}
		log.WithFields(log.Fields{"r.Password": r.Password, "r.Nbnodes": r.Nbnodes}).Debug()

		// cmd := exec.Command("/opt/slurm/launch_to_nodes.sh", r.Password, r.Nbnodes)
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		// err := cmd.Start()
		// if err != nil {
		// 	log.Error(err.Error())
		// }

		select {
		case c := <-Ch:

			myJSON, jerr := json.Marshal(c)
			if jerr != nil {
				log.Error(jerr)
			}
			Wsrw.Send([]byte(myJSON))
			if c.Result == "ok" {
				Wsrw.Close()
			}
		}
	}

}
