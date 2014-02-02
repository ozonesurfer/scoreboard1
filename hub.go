package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type hub struct {
	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection
}

var h = hub{
	broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func convert_to_bytes() ([]byte, error) {
	jstring := `{"home":"` + strconv.FormatInt(Scores[0], 10) + `","visitor":"` + strconv.FormatInt(Scores[1], 10) + `"}`
	j, err := json.Marshal(jstring)
	if err != nil {
		return nil, err
	}
	return j, nil
}
func (h *hub) run(ch chan Input) {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
			j, err := convert_to_bytes()
			if err != nil {
				fmt.Println("Json error:", err)
				os.Exit(1)
			}
			select {
			case c.send <- j:
			}
		case c := <-h.unregister:
			delete(h.connections, c)
			close(c.send)
		case i := <-ch:
			fmt.Println("Received", i)
			x, e := strconv.Atoi(i.Team)
			if e != nil || (x > len(Scores) || x < 1) {
				fmt.Println("Invalid entry")
				os.Exit(1)
			}
			x--
			scored, _ := strconv.ParseInt(i.Scored, 10, 64)
			if scored == 0 {
				fmt.Println("Game over!")
				os.Exit(0)
			}
			Scores[x] += scored
			j, err := convert_to_bytes()
			if err != nil {
				fmt.Println("Json error:", err)
				os.Exit(1)
			}
			for c := range h.connections {
				select {
				case c.send <- j:
				default:
					delete(h.connections, c)
					close(c.send)
					go c.ws.Close()

				}
			}
		}
	}
}
