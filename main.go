// concurrency3
package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Enter: TeamNumber<space>PointsJustScored
func InputFunc(c chan Input) {
	//	r := bufio.NewReader(os.Stdin)
	for {
		var a, b string
		fmt.Scanf("%s %s\n", &a, &b)
		c <- Input{Team: a, Scored: b}

	}
}

var addr = flag.String("addr", ":8080", "http service address")
var homeTempl = template.Must(template.ParseFiles("home.html"))

func homeHandler(c http.ResponseWriter, req *http.Request) {
	homeTempl.Execute(c, req.Host)
}

func OutputFunc(c chan Input) {
	for {
		select {
		case i := <-c:
			fmt.Println("Received", i)
			x, _ := strconv.Atoi(i.Team)
			x--
			scored, _ := strconv.ParseInt(i.Scored, 10, 64)
			if scored == 0 {
				fmt.Println("Bye!")
				os.Exit(0)
			}
			Scores[x] += scored
			fmt.Println("The score is now", Scores)
			time.Sleep(time.Second * 5)

		default:
		}
	}
}

var Scores [2]int64
var ch chan Input

type Input struct {
	Team   string
	Scored string
}

func main() {
	ch = make(chan Input)
	flag.Parse()
	go InputFunc(ch)
	//	go OutputFunc(ch)
	go h.run(ch)
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", wsHandler)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
	//	fmt.Println("Hello World!")
}
