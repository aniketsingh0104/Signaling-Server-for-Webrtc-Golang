package main

import (
	"net/http"
	"os"

	"WebRTCSignaling/signaling"

	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", signaling.WebSocketHandler)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8000"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})
	handler := c.Handler(mux)
	http.ListenAndServe(":"+os.Getenv("PORT"), handler)
}
