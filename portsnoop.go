package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	version = "developement version 0.1"
)

func main() {

	// CL Arguments
	host := flag.String("h", "", "Defines host either by domain or ip. Required to run program.")
	port := flag.Int("p", 0, "Defines port for lookup..")
	portSeries := flag.String("p*", "", "Defines series of ports separated with comma [,].")
	portRange := flag.String("pr", "", "Defines range of ports.")
	timeOut := flag.Int("to", 0, "Defines connection timeout in millseconds.")
	showVersion := flag.Bool("version", false, "Shows the current version of program.")
	help := flag.Bool("help", false, "Shows this help.")
	byPassLogo := flag.Bool("l", false, "Hides logo from startup.")

	flag.Parse()

	defaultTimeOut := 100 * time.Millisecond

	if *showVersion {
		fmt.Printf("PortSnoop version %s\n", version)
		os.Exit(0)
	}

	if !*byPassLogo {
		printLogo()
	}

	// Help

	if *help {
		flag.Usage()
		return
	}

	if *timeOut != 0 {
		defaultTimeOut = time.Duration(*timeOut) * time.Millisecond
	}

	fmt.Printf("Timeout set to %fs.\n", defaultTimeOut.Seconds())

	// IP address lookup for given host

	if *host != "" {

		ipAddresses, err := net.LookupHost(*host)

		if err != nil {
			fmt.Println("\x1b[31mHost error\x1b[0m:", err)
			return
		}

		for _, ip := range ipAddresses {
			fmt.Printf("The IP address of %s is %s\n", *host, ip)
		}

	} else {
		fmt.Println("Host not provided. Program can't run. Please check syntax or flag -help.")
		os.Exit(0)
	}

	// Singular port lookup lookup

	if *port != 0 && *host != "" {
		if snoopPort(*host, *port, defaultTimeOut) {
			fmt.Printf("Given port %d : \x1b[36m OPEN \x1b[39m \n", *port)
		} else {
			fmt.Printf("Given port %d : \x1b[31m CLOSED \x1b[39m \n", *port)
		}
	}

	if *portSeries != "" && *host != "" {
		fmt.Println("Scanning series of ports:\n-------------------------")
		portSplit := strings.Split(*portSeries, ",")
		for i := 0; i < len(portSplit); i++ {
			portInput, err := strconv.Atoi(portSplit[i])
			if err != nil {
				fmt.Println(portSplit[i], "is not a port. Check input.")
			} else {
				if snoopPort(*host, portInput, defaultTimeOut) {
					fmt.Printf("Given port %d : \x1b[36m OPEN \x1b[39m \n", portInput)
				} else {
					fmt.Printf("Given port %d : \x1b[31m CLOSED \x1b[39m \n", portInput)
				}
			}
		}

	}

	if *portRange != "" && *host != "" {

		// Ports by range
		portSplit := strings.Split(*portRange, "-")

		if len(portSplit) != 2 {

			fmt.Println("Invalid range format. \nFormat needs to be [p1]-[p2] (example: -pr 1-999).")

		} else {

			minPort, err1 := strconv.Atoi(portSplit[0])
			maxPort, err2 := strconv.Atoi(portSplit[1])

			if err1 != nil || err2 != nil {
				fmt.Println("Port range conversion error.\n", err1, "\n", err2)
				return
			}

			fmt.Printf("Port range set to \x1b[45m \x1b[30m %d \x1b[49m -  \x1b[45m \x1b[30m %d \x1b[49m ", minPort, maxPort)

			for port := minPort; port <= maxPort; port++ {
				if snoopPort(*host, port, defaultTimeOut) {
					fmt.Printf("\n\x1b[37m[\x1b[0mPort %d :\x1b[36m OPEN \x1b[0m\x1b[37m]\x1b[0m", port)
				}
			}
		}
	}
}

func snoopPort(host string, port int, timeout time.Duration) bool {

	addr := fmt.Sprintf("%s:%d", host, port)

	con, err := net.DialTimeout("tcp", addr, timeout)

	if err != nil {
		return false
	}

	con.Close()
	return true

}

func printLogo() {
	logoLines := []string{
		"\x1b[36m ██▓███   ▒█████   ██▀███  ▄▄▄█████▓  ██████  ███▄    █  ▒█████   ▒█████   ██▓███  ",
		"▓██░  ██▒▒██▒  ██▒▓██ ▒ ██▒▓  ██▒ ▓▒▒██    ▒  ██ ▀█   █ ▒██▒  ██▒▒██▒  ██▒▓██░  ██▒",
		"▓██░ ██▓▒▒██░  ██▒▓██ ░▄█ ▒▒ ▓██░ ▒░░ ▓██▄   ▓██  ▀█ ██▒▒██░  ██▒▒██░  ██▒▓██░ ██▓▒",
		"▒██▒ ░  ░░ ████▓▒░░██▓ ▒██▒  ▒██▒ ░ ▒██████▒▒▒██░   ▓██░░ ████▓▒░░ ████▓▒░▒██▒ ░  ░",
		"\x1b[36;1m▒▓▒░ ░  ░░ ▒░▒░▒░ ░ ▒▓ ░▒▓░  ▒ ░░   ▒ ▒▓▒ ▒ ░░ ▒░   ▒ ▒ ░ \x1b[34mBY\x1b[36;1m▒\x1b[34mKalevi\x1b[36;1m░▒░▒░ ▒▓▒░ ░  ░",
		"░▒ ░       ░ ▒ ▒░   ░▒ ░ ▒░    ░    ░ ░▒  ░ ░░ ░░   ░ ▒░  ░ ▒ ▒░   ░ ▒ ▒░ ░▒ ░ ",
		"░░       ░ ░ ░ ▒    ░░   ░   ░      ░  ░  ░     ░   ░ ░ ░ ░ ░ ▒  ░ ░ ░ ▒  ░░  \x1b[0m",
		"+---------------------------------------------------------------------------------+",
	}

	for _, line := range logoLines {
		fmt.Println(line)
	}

}
