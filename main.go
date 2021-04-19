package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/sse"
	"github.com/kardianos/service"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const version = "2021.2.1.19"
const programName = "Terminal local webservice"
const programDesription = "Display local web for rpi terminals"

type Page struct {
	Title string
	Body  []byte
}

type program struct{}

func (p *program) Start(s service.Service) error {
	LogInfo("MAIN", "Starting "+programName+" on "+s.Platform())
	go p.run()
	return nil
}

func (p *program) run() {
	CreateConfigIfNotExists()
	router := httprouter.New()
	timeStreamer := sse.New()
	networkDataStreamer := sse.New()
	router.GET("/", index)
	router.GET("/screenshot", Screenshot)
	//router.GET("/password", Password)
	//router.GET("/changenetwork", ChangeNetwork)
	//router.GET("/changenetworktodhcp", ChangeNetworkToDhcp)
	router.POST("/restart", Restart)
	router.POST("/shutdown", Shutdown)
	router.ServeFiles("/html/*filepath", http.Dir("html"))
	router.ServeFiles("/css/*filepath", http.Dir("css"))
	router.ServeFiles("/js/*filepath", http.Dir("js"))
	router.GET("/image.png", image)
	router.Handler("GET", "/listen", timeStreamer)
	router.Handler("GET", "/networkdata", networkDataStreamer)
	go StreamTime(timeStreamer)
	go StreamNetworkData(networkDataStreamer)
	LogInfo("MAIN", "Server running")
	_ = http.ListenAndServe(":9999", router)
}

func CreateConfigIfNotExists() {
	configDirectory := filepath.Join(".", "config")
	configFileName := "config.json"
	configFullPath := strings.Join([]string{configDirectory, configFileName}, "/")

	if _, checkPathError := os.Stat(configFullPath); checkPathError == nil {
		LogInfo("MAIN", "Config file already exists")
	} else if os.IsNotExist(checkPathError) {
		LogWarning("MAIN", "Config file does not exist, creating")
		mkdirError := os.MkdirAll(configDirectory, 0777)
		if mkdirError != nil {
			LogError("MAIN", "Unable to create directory for config file: "+mkdirError.Error())
		} else {
			LogInfo("MAIN", "Directory for config file created")
			data := ServerIpAddress{
				ServerIpAddress: "",
			}
			file, _ := json.MarshalIndent(data, "", "  ")
			writingError := ioutil.WriteFile(configFullPath, file, 0666)
			LogInfo("MAIN", "Writing data to JSON file")
			if writingError != nil {
				LogError("MAIN", "Unable to write data to config file: "+writingError.Error())
			} else {
				LogInfo("MAIN", "Data written to config file")
			}
		}
	} else {
		LogError("MAIN", "Config file does not exist")
	}
}

func LoadSettingsFromConfigFile() string {
	configDirectory := filepath.Join(".", "config")
	configFileName := "config.json"
	configFullPath := strings.Join([]string{configDirectory, configFileName}, "/")
	readFile, _ := ioutil.ReadFile(configFullPath)
	ConfigFile := ServerIpAddress{}
	_ = json.Unmarshal(readFile, &ConfigFile)
	ServerIpAddress := ConfigFile.ServerIpAddress
	return ServerIpAddress
}

func (p *program) Stop(s service.Service) error {
	LogInfo("MAIN", "Stopped on platform "+s.Platform())
	return nil
}

func main() {
	serviceConfig := &service.Config{
		Name:        programName,
		DisplayName: programName,
		Description: programDesription,
	}
	prg := &program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		LogError("MAIN", err.Error())
	}
	err = s.Run()
	if err != nil {
		LogError("MAIN", "Problem starting "+serviceConfig.Name)
	}
}

func StreamNetworkData(streamer *sse.Streamer) {
	timing := 20
	timeToSend := "20"
	refreshDone := true
	for {
		LogInfo("STREAM", "Streaming network data")
		start := time.Now()
		interfaceIpAddress, interfaceMask, interfaceGateway, dhcpEnabled := GetNetworkData()
		interfaceServerIpAddress := LoadSettingsFromConfigFile()
		serverAccessible, url, interfaceServerIpAddress := CheckServerIpAddress(interfaceServerIpAddress)
		if serverAccessible && !HomepageLoaded {
			timing = 0
			timeToSend = strconv.Itoa(timing)
			url = "http://localhost:9999"
		} else if serverAccessible && HomepageLoaded {
			timing--
			timeToSend = strconv.Itoa(timing)
			refreshDone = false
		} else if !HomepageLoaded {
			timing = 20
			timeToSend = strconv.Itoa(timing)
		} else if !serverAccessible {
			if !refreshDone {
				timing = 0
				url = "http://localhost:9999"
				timeToSend = strconv.Itoa(timing)
				refreshDone = true
			} else {
				timing = 20
				timeToSend = strconv.Itoa(timing)
			}
		}
		if timing < 0 {
			timing = 20
			timeToSend = strconv.Itoa(timing)
		}
		streamer.SendString("", "networkdata", interfaceIpAddress+";"+interfaceMask+";"+interfaceGateway+";"+dhcpEnabled+";"+timeToSend+";"+url+";"+interfaceServerIpAddress)
		LogInfo("STREAM", "Stream done in "+time.Since(start).String())
		time.Sleep(1 * time.Second)
	}
}

func StreamTime(streamer *sse.Streamer) {
	for {
		streamer.SendString("", "time", time.Now().Format("15:04:05"))
		time.Sleep(1 * time.Second)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles("html/" + tmpl + ".html")
	_ = t.Execute(w, p)
}
func image(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	http.ServeFile(writer, request, "image.png")
}
