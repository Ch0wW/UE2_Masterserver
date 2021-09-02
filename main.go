package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"time"
)

var ClientConnections map[*UnrealConnection]UnrealConnection

type UnrealConnection struct {
	buffer    []byte
	bufferpos int
	bufferlen int

	// incoming chan string
	conn *net.TCPConn

	Status     ConnectionStatus
	AnswerType SV_AnswerType

	protocol *UE2_GameProtocol
}

func Client_Initialize(c *net.TCPConn) (*UnrealConnection, error) {

	err := c.SetReadDeadline(time.Now().Add(1 * time.Second))

	if err != nil {
		return nil, err
	}

	client := &UnrealConnection{
		// incoming: make(chan string),
		buffer: make([]byte, 1024),
		conn:   c,

		Status: CTMS_NEEDSAUTH,

		protocol: &UE2_GameProtocol{
			protocol:   Protocol(PROTOCOL_NONE),
			clienttype: CTYPE_UNKNOWN,
		},
		//AnswerType: MSTC_UNKNOWN,
		//ClientType: CTYPE_UNKNOWN,
	}

	client.Listen()

	return client, nil
}

func (cl *UnrealConnection) Listen() {

	// Initial Sending of challenge
	err := cl.SendInitialChallenge()

	if err != nil {
		cl.CloseConnection()
	}

	// Handle the connection
	go cl.HandleConnection()
}

func (client *UnrealConnection) CloseConnection() {
	client.conn.Close()
	delete(ClientConnections, client)
	if client.conn != nil {
		client.conn = nil
	}
	client = nil

	fmt.Println("CLOSED CONNECTION")
}

func (client *UnrealConnection) HandleConnection() {

	for {
		// Stop it whenever the client disconnected.
		if client.Status == CTMS_TERMINATED {
			break
		}

		buffer := make([]byte, 1024)
		bufsize, _ := client.conn.Read(buffer)

		if bufsize <= 0 {
			//fmt.Println("server has no data to answer with")
			continue
		}

		// Skip the 4 first bytes. They have no purpose at all other
		// than indicating the size of the packet request.
		client.buffer = buffer[4:bufsize]
		client.bufferlen = bufsize
		client.bufferpos = 0

		// Printing Hex dump of bytes for debugging
		fmt.Println(hex.Dump(client.buffer))

		client.ReadMessage()
	}

	client.CloseConnection()
}

func main() {

	fmt.Println("=======================")
	fmt.Println(" Unreal Engine 2.X Masterserver")
	fmt.Println("=======================")

	BotConfig_Init()

	sAddr, err := net.ResolveTCPAddr("tcp", ":27900")
	if err != nil {
		log.Fatalln(err)
	}

	listener, _ := net.ListenTCP("tcp", sAddr)
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
		}

		// Initialize client structure
		client, err := Client_Initialize(conn)
		if err != nil {
			fmt.Println(err)
		}

		for clientList, _ := range ClientConnections {
			if clientList.conn == nil {
				client.conn = conn
				fmt.Println("Connected")
			}
		}

		fmt.Println("New size of clients:", len(ClientConnections))
	}

}
