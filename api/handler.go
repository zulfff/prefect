package api

import (
	"net/http"
	"strings"
	"time"
	"prefect/services/parser"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true
		}
		if strings.HasPrefix(origin, "http://localhost:8080") ||
			strings.HasPrefix(origin, "https://localhost:8080") {
			return true
		}
		return false
	},
}

func StreamStats(w http.ResponseWriter, r *http.Request) {
    connection, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    defer connection.Close()

    for {
        data := parser.SysDataParser()

        if err := connection.WriteJSON(data); err != nil {
            break 
        }

        time.Sleep(1 * time.Second) 
    }
}