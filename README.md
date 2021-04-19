# Terminal Local WebService Rpi

## 1. Prepare Raspberry Pi
* update and upgrade using `sudo apt-get update && sudo apt-get upgrade`
* install raspbian lite using Raspberry Pi Imager from [official site](https://www.raspberrypi.org/software/)
* using `sudo raspi-config`
  * enable console autologin
  * enable ssh
  * enable overscan (or disable underscan)
* install maim using `sudo apt-get install maim`
* install network manager using `sudo apt-get install network-manager`
* enable network manager as service using `sudo systemctl enable NetworkManager`
* start network manager as service using `sudo systemctl start NetworkManager`
* install ufw using `sudo apt-get install ufw`
* enable port 9999 using `sudo ufw allow 9999`
* reboot using `sudo reboot now`

## 2. Install Chromium in kiosk mode
* install prerequisites using `sudo apt-get install --no-install-recommends xserver-xorg x11-xserver-utils xinit openbox`
* install chromium using `sudo apt-get install --no-install-recommends chromium-browser`
* edit autostart file using `sudo nano /etc/xdg/openbox/autostart`, insert those lines:
```
# Disable any form of screen saver / screen blanking / power management
xset s off
xset s noblank
xset -dpms

# Allow quitting the X server with CTRL-ATL-Backspace
setxkbmap -option terminate:ctrl_alt_bksp

# Start Chromium in kiosk mode
sed -i 's/"exited_cleanly":false/"exited_cleanly":true/' ~/.config/chromium/'Local State'
sed -i 's/"exited_cleanly":false/"exited_cleanly":true/; s/"exit_type":"[^"]\+"/"exit_type":"Normal"/' ~/.config/chromium/Default/Preferences
chromium-browser --disable-infobars --kiosk 'http://localhost:9999'
```
* make everything start on boot using `sudo nano .bash_profile` , insert this line:
```
[[ -z $DISPLAY && $XDG_VTNR -eq 1 ]] && startx -- -nocursor
```
####TIP: by pressing `Ctrl-Alt-Backspace` you can kill chromium and get into command line

## 3. Copy program data to Raspberry
* copy files from terminal_local_webservice/rpi directory to raspberry pi /home/pi
  * rpi_linux into /home/pi
  * /html/* into /home/pi/*
  * /css/* into /home/pi/*
  * /js/* into /home/pi/*
  * example copying js directory using scp: `scp -r js pi@192.168.86.249:/home/pi`
  
## 4. Make program run as service
* create new file using `sudo nano /lib/systemd/system/zapsi.service`, insert those lines:
```
[Unit]
Description=Zapsi Service
ConditionPathExists=/home/pi/rpi_linux
After=network.target
[Service]
Type=simple
User=pi
Group=pi
LimitNOFILE=1024
Restart=on-failure
RestartSec=10
startLimitIntervalSec=60
WorkingDirectory=/home/pi
ExecStart=/home/pi/rpi_linux
PermissionsStartOnly=true
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=zapsi_service
[Install]
WantedBy=multi-user.target
```
* make sure that file is executable using `sudo chmod 755 /lib/systemd/system/zapsi.service`
* make service autostart using `sudo systemctl enable zapsi.service`
* start the service now using  `sudo systemctl start zapsi.service`
####TIP: search logs using `journalctl -f -u zapsi.service`

## 5. Clean booting screen and information
## 6. Make Raspberry Pi read-only

© 2021 Petr Jahoda
