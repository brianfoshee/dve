package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	a := firmata.NewAdaptor("/dev/tty.usbmodem144101")
	a.Connect()

	greenLED := gpio.NewLedDriver(a, "2")
	redLED := gpio.NewLedDriver(a, "4")

	greenLED.Start()
	redLED.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, html)
	})

	// path is /led/COLOR/{ON,OFF}
	http.HandleFunc("/led/", func(w http.ResponseWriter, r *http.Request) {
		segs := strings.Split(r.URL.Path, "/")

		if len(segs) != 4 {
			log.Printf("url doesn't have 3 segments %s, %v", r.URL.Path, segs)
		}

		color := segs[2]
		state := segs[3]

		switch color {
		case "red":
			if state == "on" {
				redLED.On()
			} else if state == "off" {
				redLED.Off()
			}
		case "green":
			if state == "on" {
				greenLED.On()
			} else if state == "off" {
				greenLED.Off()
			}
		}

		http.Redirect(w, r, "/", 301)
	})

	http.ListenAndServe("127.0.0.1:8080", nil)
}

var html = `
<!doctype html>
<head>
</head>
<body>
<h1>Hello Miss Kimberly's Class!</h1>
</body>
</html>
`
