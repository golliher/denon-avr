package main
// Simple program to send commands to Denon AVR and get their result

import (
	// "bufio"
	"fmt"
	"net"
	"os"
	// "time"
)

func sendCmd(cmd string, conn net.Conn)  {


	fmt.Println("Connected")
	fmt.Println("Sending: ", cmd)
	cmd = cmd + "\r"
	fmt.Fprintf(conn, cmd)


	// Bug: it's got to be an anti-pattern to have to pass connection into a routine like this.
	
	// Bug: some commands return nothing for example sending MUON when mut is already on.  Need a timeout?

	// This is a bug because there can be more than one line to read.  For example, HD?
	// How do I know when it's done sending results?


	// Commenting to prevent hands when MUON when MUON
	// status, err := bufio.NewReader(conn).ReadString('\r')
	// if err != nil {
	// 	fmt.Println("Error reading result")
	// }
	// fmt.Println("result",status)


	// return status
}

func noop() {
	fmt.Println("noop")
	return 
}


func main() {

	
	conn, err := net.Dial("tcp", "192.168.4.2:23")
	if err != nil {
		fmt.Println("Connection failed")
		os.Exit(1)
	} 
	
	
	// cmd_seq := []string{"MU?\r", "MUOFF\r", "MU?\r","MUON\r","MU?\r"}
	cmd_seq := os.Args[1:]

	for _,cmd := range cmd_seq {
		switch cmd {
		case "xboxon": { 
				sendCmd("MUOFF",conn)
				sendCmd("MV335",conn)
				sendCmd("SIGAME",conn)
				
			}
		case "xboxoff": { 
				sendCmd("MUOFF",conn)
				sendCmd("SIMPLAY",conn)
			}
		case "appletv": { 
				sendCmd("MUOFF",conn)
				sendCmd("SIMPLAY",conn)
			}

		default:
			sendCmd(cmd, conn)
		}
	}
	conn.Close()


}
