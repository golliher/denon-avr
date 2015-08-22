package main

// Simple program to send commands to Denon AVR and get their result

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var (
	gr      chan string // global receiving channel for information coming from the AVR
	conn    net.Conn    // global network connection to the AVR
	debug   = false
	verbose = true
)

func sendCmd(cmd string) {
	cmd = strings.ToUpper(cmd)
	if verbose {
		fmt.Println("Sending: ", cmd)
	}

	cmd = cmd + "\r"
	fmt.Fprintf(conn, cmd)

}

func receiver() {
	if debug {
		fmt.Println("Receiver started")
	}

	status, err := bufio.NewReader(conn).ReadString('\r')
	gr <- status
	for err == nil { // There must be more information..keep reading.
		status, err = bufio.NewReader(conn).ReadString('\r')
		gr <- status
	}

	if debug {
		fmt.Println("Receiver has stopped.")
	}

}

func printReceived() {
	if debug {
		fmt.Println("printReceived started")
	}
	for recievedMsg := range gr {
		if recievedMsg != "" {
			fmt.Println("received: ", recievedMsg)
		} else {
			if verbose {
				fmt.Println("Received no result.")
			}
		}
	}

	fmt.Println("Done printing received channel.")
}

func init() {
	if debug {
		fmt.Println("Initilizing global channels.")
	}
	gr = make(chan string)

	if debug {
		fmt.Print("Connecting..")
	}
	lconn, err := net.Dial("tcp", "192.168.4.2:23")

	if err != nil {
		fmt.Println("Connection failed")
		os.Exit(1)
	}
	if debug {
		fmt.Println("connected.")
	}
	conn = lconn // Probably a better pattern for this..

	go receiver()
	go printReceived()

}

func main() {

	// cmd_seq := []string{"MU?", "MUOFF", "MU?","MUON","MU?"}
	cmdSeq := os.Args[1:]

	defer close(gr)
	defer conn.Close()

	for _, cmd := range cmdSeq {
		switch cmd {
		case "xboxon":
			{
				sendCmd("MUOFF")
				sendCmd("MV335")
				sendCmd("SIGAME")
			}
		case "xboxoff":
			{
				sendCmd("MUOFF")
				sendCmd("SIMPLAY")
			}
		case "appletv":
			{
				sendCmd("MUOFF")
				sendCmd("SIMPLAY")
				sendCmd("Z2MPLAY")
			}
		case "ps3on":
			{
				sendCmd("MUOFF")
				sendCmd("SIBD")
			}
		case "radio":
			{
				sendCmd("MUOFF")
				sendCmd("SIHDRADIO")
			}
		case "radioall":
			{
				sendCmd("MUOFF")
				sendCmd("SIHDRADIO")
				sendCmd("MV23")
				sendCmd("Z2HDRADIO")

			}

		default:
			go sendCmd(cmd)
			// DENON API says we will have an answer after 200ms
			time.Sleep(210 * time.Millisecond)
		}

		// Do we need to wait between sending commands?
		// Probably not, but makes it easier to see whats going on during dev
		if debug {
			time.Sleep(1000 * time.Millisecond)
		}

	}

}
