package models

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// AppHandlerFunc is an HandlerFunc returning an AppError
type AppHandlerFunc func(http.ResponseWriter, *http.Request) *AppError

// AppError is the error type returned by the custom handlers
type AppError struct {
	Error   error
	Message string
	Code    int
}

// Resp is the JSON sent to the views
type Resp struct {
	Jobid  string `json:"jobid"`
	Node   string `json:"node"`
	Result string `json:"result"`
	Pass   string `json:"pass"`
}

// Req is the JSON receive by the websocket
type Req struct {
	Password string `json:"password"`
	Nbnodes  string `json:"nbnodes"`
}

// WSReaderWriter is a mutex websocket writer
type WSReaderWriter struct {
	upgrader websocket.Upgrader
	ws       *websocket.Conn
	mux      sync.Mutex
}

// Read reads the json from the websocket
func (wss *WSReaderWriter) Read() *Req {
	r := Req{}
	err := wss.ws.ReadJSON(&r)
	if err != nil {
		fmt.Println(err)
	}
	return &r
}

// Send writes the json to the websocket
func (wss *WSReaderWriter) Send(json []byte) {
	wss.mux.Lock()
	e := wss.ws.WriteMessage(websocket.TextMessage, json)
	if e != nil {
		fmt.Println(e)
	}
	wss.mux.Unlock()
}

// Send writes the json to the websocket
func (wss *WSReaderWriter) Close() {
	wss.ws.Close()
}

// WSReaderWriter constructor
func NewWSsender(w http.ResponseWriter, r *http.Request) *WSReaderWriter {
	// opening the websocket
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	wss, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic("error opening socket")
	}
	return &WSReaderWriter{ws: wss, mux: sync.Mutex{}}
}
