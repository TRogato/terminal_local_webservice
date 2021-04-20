package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func setupPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	interfaceIpAddress, interfaceMask, interfaceGateway, dhcpEnabled := GetNetworkData()
	interfaceServerIpAddress := LoadSettingsFromConfigFile()
	tmpl := template.Must(template.ParseFiles("html/setup.html"))
	data := HomepageData{
		IpAddress:       interfaceIpAddress,
		Mask:            interfaceMask,
		Gateway:         interfaceGateway,
		ServerIpAddress: interfaceServerIpAddress,
		Dhcp:            dhcpEnabled,
		DhcpChecked:     "",
		Version:         version,
	}
	if strings.Contains(dhcpEnabled, "yes") {
		data.DhcpChecked = "checked"
	}
	fmt.Println(data)
	_ = tmpl.Execute(w, data)
}

func ChangeNetwork(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_ = r.ParseForm()
	ipaddress := r.Form["ipaddress"]
	gateway := r.Form["gateway"]
	mask := r.Form["mask"]
	serveripaddress := r.Form["serveripaddress"]
	pattern := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	if pattern.MatchString(ipaddress[0]) && pattern.MatchString(gateway[0]) {
		maskNumber := GetMaskNumberFrom(mask[0])
		exec.Command("nmcli", "con", "mod", "Wired connection 1", "ipv4.method", "manual", "ipv4.addresses", ipaddress[0]+"/"+maskNumber, "ipv4.gateway", gateway[0])
		exec.Command("nmcli", "con", "up", "Wired connection 1")
	}
	if len(serveripaddress[0]) > 0 {
		configDirectory := filepath.Join(".", "config")
		configFileName := "config.json"
		configFullPath := strings.Join([]string{configDirectory, configFileName}, "/")
		data := ServerIpAddress{
			ServerIpAddress: serveripaddress[0],
		}
		file, _ := json.MarshalIndent(data, "", "  ")
		_ = ioutil.WriteFile(configFullPath, file, 0666)
	}
	_ = r.ParseForm()
	tmpl := template.Must(template.ParseFiles("html/homepage.html"))
	data := HomepageData{
		IpAddress:       "",
		Mask:            "",
		Gateway:         "",
		ServerIpAddress: "",
		Dhcp:            "",
		Version:         version,
	}
	_ = tmpl.Execute(w, data)

}

func GetMaskNumberFrom(maskNumber string) string {
	switch maskNumber {
	case "128.0.0.0":
		return "1"
	case "192.0.0.0":
		return "2"
	case "224.0.0.0":
		return "3"
	case "240.0.0.0":
		return "4"
	case "248.0.0.0":
		return "5"
	case "252.0.0.0":
		return "6"
	case "254.0.0.0":
		return "7"
	case "255.0.0.0":
		return "8"
	case "255.128.0.0":
		return "9"
	case "255.192.0.0":
		return "10"
	case "255.224.0.0":
		return "11"
	case "255.240.0.0":
		return "12"
	case "255.248.0.0":
		return "13"
	case "255.252.0.0":
		return "14"
	case "255.254.0.0":
		return "15"
	case "255.255.0.0":
		return "16"
	case "255.255.128.0":
		return "17"
	case "255.255.192.0":
		return "18"
	case "255.255.224.0":
		return "19"
	case "255.255.240.0":
		return "20"
	case "255.255.248.0":
		return "21"
	case "255.255.252.0":
		return "22"
	case "255.255.254.0":
		return "23"
	case "255.255.255.0":
		return "24"
	case "255.255.255.128":
		return "25"
	case "255.255.255.192":
		return "26"
	case "255.255.255.224":
		return "27"
	case "255.255.255.240":
		return "28"
	case "255.255.255.248":
		return "29"
	case "255.255.255.252":
		return "30"
	case "255.255.255.254":
		return "31"
	case "255.255.255.255":
		return "32"
	}
	return "0"
}

func ChangeNetworkToDhcp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	exec.Command("nmcli", "con", "mod", "Wired connection 1", "ipv4.method", "auto")
	_ = r.ParseForm()
	tmpl := template.Must(template.ParseFiles("html/homepage.html"))
	data := HomepageData{
		IpAddress:       "",
		Mask:            "",
		Gateway:         "",
		ServerIpAddress: "",
		Dhcp:            "",
		Version:         version,
	}
	_ = tmpl.Execute(w, data)
}

func CheckServerIpAddress(interfaceServerIpAddress string) bool {
	seconds := 2
	timeOut := time.Duration(seconds) * time.Second
	_, err := net.DialTimeout("tcp", interfaceServerIpAddress, timeOut)
	if err != nil {
		return false
	}
	return true
}
