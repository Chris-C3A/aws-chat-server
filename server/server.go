package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	// Deal with an error event.
	if err != nil {
		panic(err)
	}
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	// Continuously accept a network connection from the Listener
	// and add it to the channel for handling connections.
	for {
		// conn := <-conns
		conn, err := ln.Accept()
		handleError(err)

		conns <- conn
	}
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	// So long as this connection is alive:
	// Read in new messages as delimited by '\n's
	// Tidy up each message and add it to the messages channel,
	// recording which client it came from.
	for {
		msg, err := bufio.NewReader(client).ReadString('\n')
		handleError(err)

		msgs <- Message{message: fmt.Sprintf("<%d> %s", clientid, msg), sender: clientid}

	}
}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	// create listener on tcp connections on port from commandline argument
	ln, err := net.Listen("tcp", *portPtr)
	handleError(err)

	//Create a channel for connections
	conns := make(chan net.Conn)

	//Create a channel for messages
	msgs := make(chan Message)

	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)

	id := 0

	// Start accepting connections
	go acceptConns(ln, conns)

	for {
		select {
		case conn := <-conns:
			//TODO Deal with a new connection
			// - assign a client ID
			id ++
			// - add the client to the clients channel
			clients[id] = conn
			// - start to asynchronously handle messages from this client
			go handleClient(clients[id], id, msgs)

		case msg := <-msgs:
			//TODO Deal with a new message
			// Send the message to all clients that aren't the sender
			for clientid, client := range clients {
				if clientid != msg.sender {
					fmt.Fprintln(client, msg.message)
				}
			}
		}
	}
}
