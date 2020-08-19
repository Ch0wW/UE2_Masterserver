package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
	"time"
)

const (
	SERVER_PORT    = "27900"
	CONNECTIONTYPE = "tcp"
)

var ClientConnections map[*Client]int

type ClientType int

const (
	CTYPE_UNKNOWN ClientType = iota
	CTYPE_CLIENT
	CTYPE_SERVER
)

type Client struct {
	// incoming chan string
	buf        []byte
	reader     *bufio.Reader
	writer     *bufio.Writer
	conn       net.Conn
	connection *Client

	bInitialAuth bool
	bNeedsAuth   bool

	Status     ClientStatus
	ClientType ClientType
	AnswerType SV_AnswerType
}

func NewClient(connection net.Conn) *Client {
	tmpwriter := bufio.NewWriter(connection)
	tmpreader := bufio.NewReader(connection)

	// Set the buffersize to 1024
	writer := bufio.NewWriterSize(tmpwriter, 1024)
	reader := bufio.NewReaderSize(tmpreader, 1024)

	connection.SetReadDeadline(time.Now().Add(1 * time.Second))

	client := &Client{
		// incoming: make(chan string),
		buf:    make([]byte, 1024),
		conn:   connection,
		reader: reader,
		writer: writer,

		bInitialAuth: true,
		Status:       CTMS_NEEDSAUTH,
		AnswerType:   MSTC_UNKNOWN,
		ClientType:   CTYPE_UNKNOWN,
	}

	client.Listen()

	return client
}

func (client *Client) Listen() {

	// Initial Sending of challenge
	client.SendInitialChallenge()

	// Handle the connection
	go client.HandleConnection()
}

func (client *Client) CloseConnection() {
	client.conn.Close()
	delete(ClientConnections, client)
	if client.connection != nil {
		client.connection.connection = nil
	}
	client = nil

	fmt.Println("CLOSED CONNECTION")
}

func (client *Client) HandleConnection() {
	for {

		// Stop it whenever the client disconnected.
		if client.Status == CTMS_TERMINATED {
			break
		}

		buffer := make([]byte, 1024)
		bufsize, _ := client.reader.Read(buffer)

		if bufsize <= 0 {
			//fmt.Println("server has no data to answer with")
			continue
		}

		client.buf = buffer[0:bufsize]

		// Printing Hex dump of bytes for debugging
		fmt.Println(hex.Dump(client.buf))

		client.ReadChallenge(client.buf, bufsize)
		client.ReadMOTDChallenge(client.buf)
		client.Authentication_WriteData()
	}

	client.CloseConnection()
}

func main() {

	fmt.Println("=======================")
	fmt.Println(" Unreal Engine 2.X Masterserver")
	fmt.Println("=======================")

	ClientConnections = make(map[*Client]int)
	listener, _ := net.Listen("tcp", ":"+SERVER_PORT)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("ERROR: ", err.Error())
		}
		client := NewClient(conn)
		for clientList, _ := range ClientConnections {
			if clientList.connection == nil {
				client.connection = clientList
				clientList.connection = client
				fmt.Println("Connected")
			}
		}
		ClientConnections[client] = 1
		fmt.Println("New size of clients:", len(ClientConnections))
	}

}
