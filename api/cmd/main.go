package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

//// Structure

// Meta Structures
type ComputeState struct {
	ID string
	IsRunning bool
	LastActive time.Time
	Mu sync.Mutex // Lock or unlock mutual exclusivity (whether one OR more threads can access)
}

type securityConfig struct {
	api_key string
	accepted_origin string
}

type APIServer struct {
	Router *mux.Router
	ComputeState *ComputeState
	securityConfig *securityConfig
	Upgrader websocket.Upgrader
}

// Request Structures
type ControlRequest struct {
	DeviceID string `json:"device_id"` // Identify specific client machine
	Timestamp string `json:"timestamp"` // Log time
	Run bool `json:"run"`
}

type InferenceRequest struct {
	DeviceID string `json:"device_id"` // Identify specific client machine
	Timestamp string `json:"timestamp"` // Log time
	Prompt string `json:"prompt"` // Prompt that we want to respond to
}

// Response Structures
type StatusResponse struct {
	WebSocketURL string `json:"websocket_url"`
	ComputeInstance string `json:"compute_instance"`
	Status string `json:"status"`
	Ready bool `json:"ready"`
	CostPerHour float64 `json:"cost_per_hour"`
	IdleAfterMin float64 `json:"idle_after_min"`
}

type InferenceResponse struct {
	Status string `json:"status"`
	Response string `json:"response"`
	Latency string `json:"latency"`
}

//// Functionality

// Server
func LoadSecurityConfig() (*securityConfig, error){
	err := godotenv.Load(".env") 
	if err != nil {
		return nil, err
	}

	security_config := securityConfig{
		api_key: os.Getenv("API_KEY"),
		accepted_origin: os.Getenv("ACCEPTED_ORIGIN"),
	}

	return &security_config
}

func NewAPIServer() (*APIServer, error) {

	// Initialize Compute State
	compute_state := ComputeState{
		ID: "",
		IsRunning: false,
		LastActive: time.Now(),
		Mu: sync.Mutex{},
	}

	// Load and Initialize the Security Config
	security, err := LoadSecurityConfig()
	if err != nil {
		return nil, err
	}

	// Initialize Websocket Upgrader
	var upgrader = websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if origin == "" {
				return False
			}
			return (origin == security.accepted_origin)
		},
	}

	// Create the API Server
	api_server := APIServer{
		Router: mux.NewRouter(),
		ComputeState: &compute_state,
	}
	
	return &api_server, nil
}


func (api *APIServer) handleControlRequest(w http.ResponseWriter, r *http.Request) {

	var control_request ControlRequest

	if err := json.NewDecoder(r.Body).Decode(&control_request); err != {
		log.Println("control request json decoding error", err)
		http.Error(w, "invalid control request body", http.StatusBadRequest)
		return 
	}

	api.ComputeState.Mu.Lock()
	is_running := api.ComputeState.IsRunning
	api.ComputeState.Mu.Unlock()

	
	if !is_running && control_request.Run {
		//
		go s.initVastAICompute(control_request.DeviceID) // Start a concurrent thread that initializes the VastAI compute

		wsURL := fmt.Sprintf("ws://%s/status/%s", control_request.Host, control_request.DeviceID) // Create URL for websocket channel
		json.NewEncoder(w).Encode(StatusResponse{
			Status: "init",
			WebSocketURL: wsURL,
		})

		return
		//
	} else if is_running && control_request.Run {
		log.Println("trying to RUN an already RUNNING compute error")
		return 

	} else if !is_running && !control_request.Run {
		log.Println("trying to STOP an already IDLE compute error")
		return 
		
	} else if is_running && !control_request.Run {
		//
		go s.stopVastAICompute(control_request.DeviceID)
		return
		//
	}
}

func (api *APIServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader
}

func respondHandler(w http.ResponseWriter, r *http.Request) {
	var prompt PromptRequest

	if err := json.NewDecoder(r.Body).Decode(&prompt); err != nil {
		log.Println("Request Json Decoding Error: ", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Println(prompt)

	response := map[string]string{"prompt": "Prompt recieved succesfully"}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Request Json Encoding Error:", err)
		http.Error(w, "Invalid request body", http.StatusInternalServerError)
		return
	}
}

func main() {
	port := ":8000"

	api, err = NewAPIServer()
	if err != nil {
		log.Println("Starting Server Error: ", err)
	}

	api.r.HandleFunc("/control", respondHandler).Methods("POST")

	log.Printf("Server started succesfully at port: %s", port)
	log.Printf("Ready to recieve requests!")
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal("Server failed to start at port: ", port)
	}
}
