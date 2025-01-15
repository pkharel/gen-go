package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gtuk/discordwebhook"
	"gopkg.in/yaml.v3"
)

// Config file struct
type Config struct {
	Discord    string `yaml:"discord"`
	LocationID int    `yaml:"location"`
}

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	configPtr := flag.String("config", "config.yaml", "Config file")

	flag.Parse()

	configFile, err := os.ReadFile(*configPtr)
	if err != nil {
		slog.Error(err.Error())
	}

	var config Config

	err = yaml.Unmarshal(configFile, &config)

	if err != nil {
		slog.Error(err.Error())
	}

	// Get Slots for location
	c := http.Client{Timeout: time.Duration(1) * time.Second}
	req, err := http.NewRequest("GET", "https://ttp.cbp.dhs.gov/schedulerapi/slots", nil)
	if err != nil {
		slog.Error(err.Error())
	}
	q := req.URL.Query()
	q.Add("orderBy", "soonest")
	q.Add("limit", "1")
	q.Add("minimum", "1")
	q.Add("locationId", strconv.Itoa(config.LocationID))
	req.URL.RawQuery = q.Encode()

	resp, err := c.Do(req)

	if err != nil {
		slog.Error(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error(err.Error())
	}

	var slots []Slot
	json.Unmarshal(body, &slots)

	if len(slots) == 0 {
		slog.Info("No Appointments Available")
	} else {
		slog.Info("Appointments available!")
	}

	username := "BotUser"
	content := "Appointments Available!"
	url := config.Discord

	message := discordwebhook.Message{
		Username: &username,
		Content:  &content,
	}

	err = discordwebhook.SendMessage(url, message)
	if err != nil {
		log.Fatal(err)
	}
}

type Slot struct {
	LocationID     int
	StartTimestamp string
	EndTimestamp   string
	Active         bool
	Duration       int
	RemoteInd      bool
}
