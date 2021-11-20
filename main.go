package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"time"
)

var ServerList *[]UnrealServerData

type UnrealServerData struct {
	IP               *net.IP
	Port             uint16
	QueryPort        uint16
	GameSpyQueryPort uint16
}

var ConnectionList *[]UnrealConnection

var debugmode = true

type ServerUDPRequest struct {
	Code             uint32
	GamePort         uint16
	GameQueryPort    uint16
	GameSpyQueryPort uint16
}

type UnrealConnection struct {
	buffer    []byte
	bufferpos int
	bufferlen int

	UT2K4_userVerified bool

	// incoming chan string
	conn    *net.TCPConn
	udpconn *net.UDPConn

	Status     ConnectionStatus
	AnswerType SV_AnswerType

	Protocol *UE2_GameProtocol

	SV_UDPInfo *ServerUDPRequest
}

func Client_HandleConnection(co *net.TCPConn) {

	client := &UnrealConnection{
		conn:   co,
		Status: CTMS_NEEDSAUTH,

		Protocol: &UE2_GameProtocol{
			protocol:   Protocol(PROTOCOL_NONE),
			clienttype: CTYPE_UNKNOWN,
			language:   "int",
		},
	}

	// Server always does the first step...
	err := client.SendInitialChallenge()
	if err != nil {
		co.Close()
		return
	}

	// Then, loop our connection...
	for {
		err := co.SetReadDeadline(time.Now().Add(60 * time.Second))

		if err != nil {
			break
		}

		defsize := 4

		buffer := make([]byte, 1024)
		size, err := co.Read(buffer)
		if err != nil {
			break
		}

		newbuf := buffer[defsize:size]

		if debugmode {
			//fmt.Println(size, newbuf)
		}

		client.buffer = newbuf
		client.bufferlen = size - defsize
		client.bufferpos = 0

		fmt.Println(hex.Dump(newbuf))

		// Process the message and check if something goes wrong...
		err = client.ReadMessage()
		if err != nil {
			break
		}
	}

	// Safely close the connection
	fmt.Println("Connection closed...")
	co.Close()
}

func main() {

	fmt.Println("=======================")
	fmt.Println(" Unreal Engine 2.X Masterserver")
	fmt.Println(" Version 0.1a by Ch0wW")
	fmt.Println("=======================")

	go UDP_main()

	BotConfig_Init()

	sAddr, err := net.ResolveTCPAddr("tcp", ":27900")
	if err != nil {
		log.Fatalln(err)
	}

	listener, _ := net.ListenTCP("tcp", sAddr)
	defer listener.Close()
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			conn.Close()
		}

		fmt.Println("Client connected")
		go Client_HandleConnection(conn)
		// Initialize client structure
	}
}

func UDP_main() {
	sAddr, err := net.ResolveUDPAddr("udp", ":27900")
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := net.ListenUDP("udp", sAddr)
	conn.SetReadBuffer(1048576)
	if err != nil {
		// handle error
	}
	defer conn.Close()

	var buf [1024]byte
	for {
		rlen, remote, err := conn.ReadFromUDP(buf[:])

		if err != nil {
			panic(err)
		}

		if debugmode {
			fmt.Println(rlen, remote, buf[:rlen])
		}

		// Do stuff with the read bytes
		go UDP_HandleConnection(conn, remote, buf, rlen)
	}

}

func UDP_HandleConnection(co *net.UDPConn, addr *net.UDPAddr, buffer [1024]byte, rlen int) {

	c := &UnrealConnection{
		udpconn: co,
	}

	_, err := c.ReadLong()
	if err != nil {
		fmt.Println("Unknown code read...")
		return
	}

	porttype, err := c.ReadByte()
	if err != nil {
		fmt.Println("Unable to read port type #", porttype)
		return
	}

	code, err := c.ReadLong()
	if err != nil {
		fmt.Println("Unable to read code", porttype)
		return
	}

}
