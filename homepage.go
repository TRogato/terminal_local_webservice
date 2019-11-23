package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/kbinani/screenshot"
	"html/template"
	"image/png"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type ServerIpAddress struct {
	ServerIpAddress string
}

type HomepageData struct {
	IpAddress       string
	Mask            string
	Gateway         string
	ServerIpAddress string
	Timer           string
	Url             string
	Dhcp            string
}

func Screenshot(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	LogInfo("MAIN", "Generating screenshot")
	n := screenshot.NumActiveDisplays()
	LogInfo("MAIN", "Displays: "+strconv.Itoa(n))

	for i := 0; i < n; i++ {
		img, err := screenshot.CaptureDisplay(i)
		if err != nil {
			LogError("MAIN", "Error generating screenshot: "+err.Error())
			continue
		}
		fileName := "image.png"
		file, _ := os.Create(fileName)
		defer file.Close()
		_ = png.Encode(file, img)
		LogInfo("MAIN", "Generated screenshot: "+fileName)
	}
	LogInfo("MAIN", "Generating finished")
	renderTemplate(w, "screenshot", &Page{})
}

func Restart(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	data, err := exec.Command("Powershell.exe", "Restart-Computer").Output()

	if err != nil {
		fmt.Println("Error: ", err)
	}
	LogInfo("MAIN", string(data))
}

func Shutdown(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	data, err := exec.Command("Powershell.exe", "Stop-Computer").Output()

	if err != nil {
		fmt.Println("Error: ", err)
	}
	LogInfo("MAIN", string(data))
}

func Homepage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	_ = r.ParseForm()
	tmpl := template.Must(template.ParseFiles("html/homepage.html"))

	interfaceServerIpAddress := LoadSettingsFromConfigFile()
	timer := "86400"
	url := ""

	hostName := interfaceServerIpAddress
	portNum := "80"
	seconds := 2
	timeOut := time.Duration(seconds) * time.Second

	_, err := net.DialTimeout("tcp", hostName+":"+portNum, timeOut)

	if err != nil {
		LogError("MAIN", interfaceServerIpAddress+" not accessible: "+err.Error())
		interfaceServerIpAddress += " not accessible"
	} else {
		LogInfo("MAIN", interfaceServerIpAddress+" accessible")
		timer = "20"
		url = "http://" + interfaceServerIpAddress + "/"
	}

	data := HomepageData{
		IpAddress:       "",
		Mask:            "",
		Gateway:         "",
		ServerIpAddress: interfaceServerIpAddress,
		Timer:           timer,
		Url:             url,
		Dhcp:            "",
	}
	_ = tmpl.Execute(w, data)
}

func Setup(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	_ = r.ParseForm()
	tmpl := template.Must(template.ParseFiles("html/setup.html"))

	interfaceServerIpAddress := LoadSettingsFromConfigFile()
	timer := "86400"
	url := ""

	hostName := interfaceServerIpAddress
	portNum := "80"
	seconds := 2
	timeOut := time.Duration(seconds) * time.Second

	_, err := net.DialTimeout("tcp", hostName+":"+portNum, timeOut)

	if err != nil {
		LogError("MAIN", interfaceServerIpAddress+" not accessible: "+err.Error())
		interfaceServerIpAddress += " not accessible"
	} else {
		LogInfo("MAIN", interfaceServerIpAddress+" accessible")
		timer = "20"
		url = "http://" + interfaceServerIpAddress + "/"
	}
	data := HomepageData{
		IpAddress:       "",
		Mask:            "",
		Gateway:         "",
		ServerIpAddress: interfaceServerIpAddress,
		Timer:           timer,
		Url:             url,
		Dhcp:            "",
	}
	_ = tmpl.Execute(w, data)
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
