package server

import (
	"github.com/Battle-Bunker/cyphid-snake/agent"
	"github.com/BattlesnakeOfficial/rules/client"
	"encoding/json"
	"log"
	"net/http"
	// "io"
	// "bytes"
	"os"
	"io"
)

type Server struct {
	agent *agent.SnakeAgent
}

func NewServer(agent *agent.SnakeAgent) *Server {
	return &Server{agent: agent}
}

// Middleware

const ServerID = "battlesnake/github/starter-snake-go"

func withServerID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", ServerID)
		next(w, r)
	}
}

// Start Battlesnake Server
func (s *Server) Start() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8000"
	}

	http.HandleFunc("/", withServerID(s.handleIndex))
	http.HandleFunc("/start", withServerID(s.handleStart))
	http.HandleFunc("/move", withServerID(s.handleMove))
	http.HandleFunc("/end", withServerID(s.handleEnd))

	log.Printf("Running Battlesnake at http://0.0.0.0:%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}


func (s *Server) handleStart(w http.ResponseWriter, r *http.Request) {
	log.Println("START")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleMove(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received move request with Content-Length: %d", r.ContentLength)

	if r.Body == nil {
		log.Printf("Error: request body is nil")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "empty request body"})
		return
	}

	var request client.SnakeRequest
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to read request body"})
		return
	}
	defer r.Body.Close()

	if len(bodyBytes) == 0 {
		log.Printf("Error: empty request body")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "empty request body"})
		return
	}

	if err := json.Unmarshal(bodyBytes, &request); err != nil {
		log.Printf("Error decoding move request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "unable to decode request"})
		return
	}

	var gameSnapshot agent.GameSnapshot
	if gameSnapshot = agent.NewGameSnapshot(&request); gameSnapshot == nil {
		log.Printf("Error creating game snapshot")
				w.WriteHeader(http.StatusInternalServerError)
				response := map[string]string{"error": "unable to create game snapshot"}
				json.NewEncoder(w).Encode(response)
		return
	}

	moveResponse := s.agent.ChooseMove(gameSnapshot)
	log.Printf("Turn %d: Move %s, Shout '%s'", request.Turn, moveResponse.Move, moveResponse.Shout)
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(moveResponse); err != nil {
				log.Printf("Error encoding move response: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
		}
}

func (s *Server) handleEnd(w http.ResponseWriter, r *http.Request) {
	log.Println("END")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) { 
	metadata := s.agent.Metadata

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(metadata); err != nil {
		log.Printf("Error encoding info response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}