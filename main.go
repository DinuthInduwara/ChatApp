package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var websocketUpgrader = websocket.Upgrader{
	// Apply the Origin Checker
	CheckOrigin:     checkOrigin,
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func checkOrigin(r *http.Request) bool {
	return true // Allow all origins for simplicity
	//origin := r.Header.Get("Origin")
	// switch origin {
	// case "http://127.0.0.1:" + strconv.Itoa(PORT):
	// 	return true
	// default:
	// 	return false
	// }
}
func checkPort(port int) (int, error) {
	addr := net.JoinHostPort("0.0.0.0", strconv.Itoa(port))
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		// Port is in use, try a random port
		for retries := 0; retries < 10; retries++ {
			newPort := 1024 + rand.Intn(65535-1024)
			addr = net.JoinHostPort("0.0.0.0", strconv.Itoa(newPort))
			listener, err = net.Listen("tcp", addr)
			if err == nil {
				listener.Close()
				return newPort, nil
			}
		}
		return 0, fmt.Errorf("could not find a free port after 10 retries")
	}
	listener.Close()
	return port, nil
}

var Users = make(map[string]*client)
var PORT int = 8080

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

	mux.Handle("/chat/", http.StripPrefix("/chat", http.FileServer(http.Dir("./public"))))
	mux.Handle("/", handler)

	PORT, err := checkPort(PORT)
	if err != nil {
		log.Fatal("Error finding a free port: ", err)
	}

	server := &http.Server{
		Addr:    ":" + strconv.Itoa(PORT),
		Handler: mux,
	}
	log.Printf("Server started on port %d", PORT)
	log.Printf("Open http://localhost:%d/chat/ in your browser", PORT)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
