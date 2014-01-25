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

func (h *hub) run(ch chan Input) {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			delete(h.connections, c)
			close(c.send)
			/*	case m := <-h.broadcast:
				for c := range h.connections {
					select {
					case c.send <- m:
					default:
						delete(h.connections, c)
						close(c.send)
						go c.ws.Close()
					}
				} */
		case i := <-ch:
			fmt.Println("Received", i)
			x, _ := strconv.Atoi(i.Team)
			x--
			scored, _ := strconv.ParseInt(i.Scored, 10, 64)
			Scores[x] += scored
			jstring := `{"home":"` + strconv.FormatInt(Scores[0], 10) + `","visitor":"` + strconv.FormatInt(Scores[1], 10) + `"}`
			j, err := json.Marshal(jstring)
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
