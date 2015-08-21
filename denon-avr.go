package main

// Simple program to send commands to Denon AVR and get their result

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

var (
	gr      chan string // global receiving channel for information coming from the AVR
	conn    net.Conn    // global network connection to the AVR
	debug   = false
	verbose = false
)

func sendCmd(cmd string) {
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

}

func init() {
	if debug {
		fmt.Print("Initilizing global channels.")
	}
	gr = make(chan string)

	if debug {
		fmt.Print("Connecting..")
	}
	lconn, err := net.Dial("tcp", "192.168.4.2:23")
	lconn.SetDeadline(time.Now().Add(300 * time.Millisecond)) // API spec says results should take no more than 200ms
	if err != nil {
		fmt.Println("Connection failed")
		os.Exit(1)
	}
	if debug {
		fmt.Println("connected.")
	}
	conn = lconn // Probably a better pattern for this..

	go receiver()

}

func main() {

	// cmd_seq := []string{"MU?", "MUOFF", "MU?","MUON","MU?"}
	cmdSeq := os.Args[1:]

	defer close(gr)

	// BUG: Only the first command in the cmd_seq actually works
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
			}

		default:
			sendCmd(cmd)
			printReceived()

		}
	}

	conn.Close()

}
