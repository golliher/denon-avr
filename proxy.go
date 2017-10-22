package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var (
	allClients map[*Client]int
	gr         chan string // global receiving channel for information coming from the AVR
	avr_conn   net.Conn    // global network connection to the AVR
	debug      = true
	verbose    = true
)

func sendCmd(cmd string) {

	cmd = strings.ToUpper(cmd)
	if verbose {
		fmt.Println("Sending: ", cmd)
	}

	if cmd[len(cmd)-1:] != "\r" {
		cmd = cmd + "\r"
		fmt.Println("Adding slash r to command")
	}

	fmt.Fprintf(avr_conn, cmd)

}

func receiver() {
	if debug {
		fmt.Println("Receiver started")
	}

	status, err := bufio.NewReader(avr_conn).ReadString('\r')
	gr <- status
	for err == nil { // There must be more information..keep reading.
		status, err = bufio.NewReader(avr_conn).ReadString('\r')
		gr <- status
	}

	if debug {
		fmt.Println("Receiver has stopped.")
	}

}

func init() {
	if debug {
		fmt.Println("Initilizing global channels.")
	}
	gr = make(chan string)

	if debug {
		fmt.Print("Connecting to AVR..")
	}
	lconn, err := net.Dial("tcp", "192.168.4.2:23")

	if err != nil {
		fmt.Println("Connection failed")
		os.Exit(1)
	}
	if debug {
		fmt.Println("connected.")
	}
	avr_conn = lconn // Probably a better pattern for this..

	go receiver()
}

func printReceived() {
	if debug {
		fmt.Println("printReceived started")
	}
	for recievedMsg := range gr {
		if recievedMsg != "" {
			fmt.Println("received from AVR: ", recievedMsg)
			for clientList, _ := range allClients {
				clientList.outgoing <- recievedMsg
			}

		} else {
			if verbose {
				fmt.Println("Received no result.")
			}
		}
	}

	fmt.Println("Done printing received channel.")

}

/////

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
			fmt.Println(line)
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
				fmt.Println("Connected")
			}
		}
		allClients[client] = 1
		fmt.Println(len(allClients))
	}
}
