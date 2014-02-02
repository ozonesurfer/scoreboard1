This website experiment uses Go to combine command line input with a website, using WebSockets. It simulates a scoreboard.

# You Will Need To:

Install: [Go/golang](https://code.google.com/p/go/downloads/list)

Run: <b>go get github.com/gorilla/websocket</b>

# Building

To build the executable, execute the following:

<b>go build</b>

# Running

First, run the <b>scoreboard1</b> command. Then, enter pairs of numbers; the first is either 1 for the home team or 2 for the visiting team, the second is the points just scored. Enter a zero (0) for the points just scored to exit the program.

The scoreboard is viewable on your browser at:

<i>http://localhost:8080</i>