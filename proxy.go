package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var (
	allClients map[*Client]int
	gr         chan string // global receiving channel for information coming from the AVR
	avr_conn   net.Conn    // global network connection to the AVR
	debug      = false
	verbose    = false
)

func print_verbose(m ...interface{}) {
	if verbose {
		fmt.Println(m...)
	}
}

func print_debug(m ...interface{}) {
	if debug {
		fmt.Println(m...)
	}
}

func sendCmd(cmd string) {
	cmd = strings.ToUpper(cmd)
	print_verbose("   Sending to AVR: ", cmd)

	if cmd[len(cmd)-1:] != "\r" {
		cmd = cmd + "\r"
	}

	fmt.Fprintf(avr_conn, cmd)
}

func receiver() {
	print_verbose("Receiver started")

	status, err := bufio.NewReader(avr_conn).ReadString('\r')
	gr <- status
	for err == nil { // There must be more information..keep reading.
		status, err = bufio.NewReader(avr_conn).ReadString('\r')
		gr <- status
	}
	print_verbose("Receiver has stopped.")
}

func init() {

	flag.BoolVar(&verbose, "v", false, "Enable verbose output")
	flag.Parse()

	print_debug("Initilizing global channels.")
	gr = make(chan string)

	print_verbose("Connecting to AVR.")
	lconn, err := net.Dial("tcp", "192.168.4.2:23")
	if err != nil {
		fmt.Println("Connection to AVR failed. Exiting.")
		os.Exit(1)
	}
	print_verbose("...connected to AVR.")
	avr_conn = lconn
	go receiver()
}

func printReceived() {

	print_debug("printReceived started")

	for recievedMsg := range gr {
		if recievedMsg != "" {
			print_verbose("Received from AVR: ", recievedMsg)
			for clientList, _ := range allClients {
				clientList.outgoing <- recievedMsg
			}

		} else {
			print_verbose("Received no result.")
		}
	}
	print_debug("Done printing received channel.")

}

type Client struct {
	// incoming chan string
	outgoing   chan string
	reader     *bufio.Reader
	writer     *bufio.Writer
	conn       net.Conn
	connection *Client
}

func (client *Client) Read() {
	for {
		line, err := client.reader.ReadString('\r')
		if err == nil {
			if client.connection != nil {
				client.connection.outgoing <- line
			}
			sendCmd(line)
		} else {
			break
		}
	}

	client.conn.Close()
	delete(allClients, client)
	if client.connection != nil {
		client.connection.connection = nil
	}
	client = nil
}

func (client *Client) Write() {
	for data := range client.outgoing {
		client.writer.WriteString(data)
		client.writer.Flush()
	}
}

func (client *Client) Listen() {
	go client.Read()
	go client.Write()
}

func NewClient(connection net.Conn) *Client {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	client := &Client{
		// incoming: make(chan string),
		outgoing: make(chan string),
		conn:     connection,
		reader:   reader,
		writer:   writer,
	}
	client.Listen()

	return client
}

func main() {
	allClients = make(map[*Client]int)

	go printReceived()

	fmt.Println("Starting server")

	listener, _ := net.Listen("tcp", ":23")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}
		client := NewClient(conn)
		for clientList, _ := range allClients {
			if clientList.connection == nil {
				client.connection = clientList
				clientList.connection = client
				print_verbose("CLient connected")
			}
		}
		allClients[client] = 1
		print_debug(len(allClients))
	}
}
