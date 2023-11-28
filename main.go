package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var websocketUpgrader = websocket.Upgrader{
	// Apply the Origin Checker
	CheckOrigin:     checkOrigin,
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	switch origin {
	case "http://127.0.0.1:8080":
		return true
	default:
		return false
	}
}

var Users = make(map[string]*client)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			w.WriteHeader(http.StatusNonAuthoritativeInfo)
			w.Write([]byte("No 'name' found in params"))
			return
		}
		conn, err := websocketUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		if _, ok := Users[name]; !ok {
			cli := newClient(conn, name)
			Users[name] = cli
			go cli.Read()
			if data, err := json.Marshal(&Event{Type: EventNewUser, Name: name}); err == nil {
				broadcastMessage(&data)
				log.Println(string(data), name)
			} else {
				log.Println(err)
			}
		}
	})
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
	})

	mux.Handle("/fs/", http.StripPrefix("/fs", http.FileServer(http.Dir("./public"))))
	mux.Handle("/", handler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
