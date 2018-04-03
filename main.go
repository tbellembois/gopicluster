package main

//go:generate rice embed-go

import (
	"flag"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"

	"github.com/justinas/alice"
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gopicluster/handlers"
	"github.com/tbellembois/gopicluster/models"
)

func main() {

	var (
		err              error
		logf             *os.File
		tmplBaseString   string
		tmplHeadString   string
		tmplFootString   string
		tmplHeaderString string
		tmplFooterString string
		tmplMenuString   string
		tmplHomeString   string
	)

	// getting the program parameters
	port := flag.String("port", "8080", "the port to listen")
	logfile := flag.String("logfile", "", "log to the given file")
	address := flag.String("address", "localhost", "the application address")
	debug := flag.Bool("debug", false, "debug (verbose log), default is error")
	dryrun := flag.Bool("dryrun", false, "do not call system commands")
	flag.Parse()

	// logging to file if logfile parameter specified
	if *logfile != "" {
		if logf, err = os.OpenFile(*logfile, os.O_WRONLY|os.O_CREATE, 0755); err != nil {
			log.Panic(err)
		} else {
			log.SetOutput(logf)
		}
	}

	// setting the log level
	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.ErrorLevel)
	}

	// environment creation
	env := handlers.Env{
		Templates: make(map[string]*template.Template),
		Address:   *address,
		Port:      *port,
		DryRun:    *dryrun,
	}

	// router definition
	r := mux.NewRouter()
	commonChain := alice.New(env.HeadersMiddleware, env.LogingMiddleware)
	r.Handle("/", commonChain.Then(env.AppMiddleware(env.HomeHandler))).Methods("GET")

	r.Handle("/job/start", commonChain.Then(env.AppMiddleware(env.StartHandler))).Methods("GET")
	r.Handle("/job/stop", commonChain.Then(env.AppMiddleware(env.StopHandler))).Methods("GET")
	r.Handle("/test", commonChain.Then(env.AppMiddleware(env.TestHandler))).Methods("GET")
	r.Handle("/crack", commonChain.Then(env.AppMiddleware(env.CrackHandler))).Methods("PUT")
	r.HandleFunc("/socket/", env.SocketHandler)

	// rice boxes
	cssBox := rice.MustFindBox("static/css")
	cssFileServer := http.StripPrefix("/css/", http.FileServer(cssBox.HTTPBox()))
	http.Handle("/css/", cssFileServer)

	jsBox := rice.MustFindBox("static/js")
	jsFileServer := http.StripPrefix("/js/", http.FileServer(jsBox.HTTPBox()))
	http.Handle("/js/", jsFileServer)

	imgBox := rice.MustFindBox("static/img")
	imgFileServer := http.StripPrefix("/img/", http.FileServer(imgBox.HTTPBox()))
	http.Handle("/img/", imgFileServer)

	tmplBox := rice.MustFindBox("static/templates")
	tmplHomeBox := rice.MustFindBox("static/templates/home")
	if tmplBaseString, err = tmplBox.String("base.html"); err != nil {
		log.Fatal("base.html :" + err.Error())
	}
	if tmplHeadString, err = tmplBox.String("head.html"); err != nil {
		log.Fatal("head.html :" + err.Error())
	}
	if tmplFootString, err = tmplBox.String("foot.html"); err != nil {
		log.Fatal("foot.html :" + err.Error())
	}
	if tmplHeaderString, err = tmplBox.String("header.html"); err != nil {
		log.Fatal("header.html :" + err.Error())
	}
	if tmplFooterString, err = tmplBox.String("footer.html"); err != nil {
		log.Fatal("footer.html :" + err.Error())
	}
	if tmplMenuString, err = tmplBox.String("menu.html"); err != nil {
		log.Fatal("menu.html :" + err.Error())
	}
	if tmplHomeString, err = tmplHomeBox.String("index.html"); err != nil {
		log.Fatal("home index.html :" + err.Error())
	}

	http.Handle("/", r)

	// templates compilation
	tplbase := []string{tmplBaseString,
		tmplHeadString,
		tmplFootString,
		tmplHeaderString,
		tmplFooterString,
		tmplMenuString,
	}

	tplhome := append(tplbase, tmplHomeString)
	if env.Templates["home"], err = template.New("home").Parse(strings.Join(tplhome, "")); err != nil {
		panic(err)
	}

	// initialize the channel
	handlers.Ch = make(chan models.Resp, 600)
	handlers.T = time.Now()

	if err = http.ListenAndServe(":"+*port, nil); err != nil {
		panic(err)
	}
}
