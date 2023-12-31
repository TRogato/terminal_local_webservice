# Display WebService Changelog

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/).

Please note, that this project, while following numbering syntax, it DOES NOT
adhere to [Semantic Versioning](http://semver.org/spec/v2.0.0.html) rules.

## Types of changes

* ```Added``` for new features.
* ```Changed``` for changes in existing functionality.
* ```Deprecated``` for soon-to-be removed features.
* ```Removed``` for now removed features.
* ```Fixed``` for any bug fixes.
* ```Security``` in case of vulnerabilities.


## [2022.4.3.3] - 2022-12-03

### Added
- license

## [2022.1.2.14] - 2022-01-14

### Added
- shutdown button
- demobutton with 10 additional pages

## [2022.1.2.10] - 2022-01-10

### Changed
- updated for latest RaspberryPiOS

## [2021.2.2.27] - 2021-05-27

### Changed
- updated README.md

## [2021.2.2.10] - 2021-05-10

### Changed
- removed buttons for restart and shutdown
- added page for remote network change (without on screen keyboard)
- added information about MAC
- change from port 9999 to port 80
- password entering hidden

## [2021.2.2.4] - 2021-05-04

### Changed
- checking server address
- index and setup page design

### Added
- /server function for setting only server address
- editing input fields via keyboard from cursor position

## [2021.2.2.3] - 2021-05-03

### Changed
- updated to latest go
- updated to latest go libraries

## [2021.2.1.30] - 2020-04-30

### Changed
- password field is blinking
- network is restarted after DHCP is set
- www.zapsi.eu:80 is set as default server

## [2021.2.1.29] - 2020-04-29

### Changed
- updated to check available ethernet connections on boot

## [2021.2.1.28] - 2020-04-28

### Changed
- updated to work with multiple ethernet connections

## [2021.2.1.23] - 2020-04-23

### Added
- better information about device and server network accessibility

### Changed
- code reformat and cleanup

## [2021.2.1.22] - 2020-04-22

### Changed
- streaming data only on main screen
- saving configuration

### Added
- loading and setting rpi configuration after start
- using initramfs and /ro
- info panel

## [2021.2.1.21] - 2020-04-21

### Added
- complete setup page with setting terminal to dhcp or static
- complete loading server page when available
- updated readme.md

## [2021.2.1.20] - 2020-04-20

### Added
- setup page with keyboard and proper functionality
- loading data for dhcp when cable plugged or unplugged


## [2021.2.1.19] - 2020-04-19

### Changed
- new main screen

### Added
- manual for creating the rpi image

## [2021.2.1.12] - 2020-04-12

### Changed
- rpi only

## [2020.1.3.31] - 2020-03-31

### Added
- udpated darcula.css, added "overscroll-behavior-x: none;" to body, preventing touch back

## [2020.1.3.7] - 2020-03-07

### Added
- deleting old log files

### Changed
- code reformat, pieces of code moved to proper files

### Fixed
- properly generating screenshot, not as 'root', but as a 'zapsi' user



## [2020.1.3.6] - 2020-03-06

### Added
- changelog.md
- fully working on solus linux with AsusPro
- added linux installation guide
- updated logging