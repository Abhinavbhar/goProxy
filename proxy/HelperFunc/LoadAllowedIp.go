package helperfunc

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type AllowedIp struct {
	Ip      []string `json:"ips"`
	Success bool     `json:"success"`
}

func LoadAllowedIp() []string {
	resp, err := http.Get("http://localhost:3000/allowedips")
	if err != nil {
		log.Fatal("failed to load allowed Ip (Backend server failed)", err)
	}
	defer resp.Body.Close()
	body, errReading := io.ReadAll(resp.Body)
	if errReading != nil {
		log.Fatal("error reading the ips")
	}
	var Ips AllowedIp
	errUnmarshaling := json.Unmarshal(body, &Ips)
	if errUnmarshaling != nil {
		fmt.Println("failed to unmarshal the json", errUnmarshaling)
	}
	if !Ips.Success {
		fmt.Println("Apis failed to load Ips")
	}
	return Ips.Ip
}
