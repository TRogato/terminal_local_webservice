package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/sse"
	"html/template"
	"net/http"
	"time"
)

var Interface = "Ethernet"

const deleteLogsAfter = 240 * time.Hour

type Page struct {
	Title string
	Body  []byte
}

func main() {
	LogDirectoryFileCheck("MAIN")
	router := httprouter.New()
	streamer := sse.New()
	router.GET("/", Homepage)
	router.GET("/screenshot", Screenshot)
	router.GET("/restart", Restart)
	router.GET("/setup", Setup)
	router.GET("/darcula.css", darcula)
	router.GET("/metro.min.js", metrojs)
	router.GET("/metro-all.css", metrocss)
	router.GET("/image.png", image)
	router.Handler("GET", "/listen", streamer)
	go StreamTime(streamer)
	LogInfo("MAIN", "Server running")
	_ = http.ListenAndServe(":8000", router)
}

func StreamTime(streamer *sse.Streamer) {
	for {
		streamer.SendString("", "time", time.Now().Format("15:04:05"))
		time.Sleep(1 * time.Second)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	_ = t.Execute(w, p)
}

func darcula(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	http.ServeFile(writer, request, "darcula.css")
}

func metrojs(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	http.ServeFile(writer, request, "metro.min.js")
}
func metrocss(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	http.ServeFile(writer, request, "metro-all.css")
}

func image(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	http.ServeFile(writer, request, "image.png")
}