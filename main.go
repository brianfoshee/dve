package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

// https://godoc.org/gobot.io/x/gobot/platforms/firmata
// https://godoc.org/gobot.io/x/gobot/drivers/gpio#LedDriver
func main() {
	a := firmata.NewAdaptor("/dev/tty.usbmodem144101")
	a.Connect()

	greenLED := gpio.NewLedDriver(a, "2")
	redLED := gpio.NewLedDriver(a, "4")

	greenLED.Start()
	redLED.Start()

	var homeTempl = template.Must(template.New("home").Parse(homePage))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		greenState := greenLED.State()
		redState := redLED.State()

		gs, rs, ngs, nrs := "Off", "Off", "On", "On"
		if greenState {
			gs = "On"
			ngs = "Off"
		}
		if redState {
			rs = "On"
			nrs = "Off"
		}

		data := struct {
			GreenState    string
			NotGreenState string
			RedState      string
			NotRedState   string
		}{
			GreenState:    gs,
			NotGreenState: ngs,
			RedState:      rs,
			NotRedState:   nrs,
		}

		err := homeTempl.Execute(w, data)
		if err != nil {
			log.Println("homeTempl execution error: ", err)
		}
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
			if state == "On" {
				redLED.On()
			} else if state == "Off" {
				redLED.Off()
			}
		case "green":
			if state == "On" {
				greenLED.On()
			} else if state == "Off" {
				greenLED.Off()
			}
		}

		http.Redirect(w, r, "/", 301)
	})

	http.ListenAndServe("127.0.0.1:8080", nil)
}

// TODO add forms with buttons to turn red and green on/off
var homePage = `
<!doctype html>
<head>
<style>
#container {
	display: grid;
	grid-template-columns: 400px 400px;
}
</style>
</head>

<body>

<h1>Hello, Miss Kimberly's Class!</h1>

<div id="container">
	<div id="red">
		<h2>The Red LED is {{.RedState}}</h2>
		<form action="/led/red/{{.NotRedState}}" method="post">
		<div class="button">
		  <button type="submit">Turn the Red LED {{.NotRedState}}</button>
		</div>
		</form>
	</div>

	<div id="green">
		<h2>The Green LED is {{.GreenState}}</h2>
		<form action="/led/green/{{.NotGreenState}}" method="post">
		<div class="button">
		  <button type="submit">Turn the Green LED {{.NotGreenState}}</button>
		</div>
		</form>
	</div>
</div>

</body>
</html>
`
