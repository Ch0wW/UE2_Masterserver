package main

import (
	"encoding/hex"
	"fmt"
)

const (
	AUTHINFO      = "\x03\x00\x00\x00\x02\x30\x00"                             // Unknown about why this part of code...
	AUTHDENIED    = "\x08\x00\x00\x00\x09\x44\x45\x4e\x49\x45\x44\x00"         // DENIED message
	AUTH_APPROVED = "\x0a\x00\x00\x00\x09\x41\x50\x50\x52\x4f\x56\x45\x44\x00" // APPROVED message
)

type ClientStatus int

const (
	CTMS_NEEDSAUTH ClientStatus = iota
	CTMS_WAITING
	CTMS_TERMINATED
	CTMS_LOGGED
)

type SV_AnswerType int

const (
	MSTC_UNKNOWN SV_AnswerType = -1
	MSTC_DENIED  SV_AnswerType = iota
	MSTC_APPROVED
	MSTC_UPGRADE // Not Set
	MSTC_MSLIST  // Not Set
	MSTC_MOTD
)

func (client *Client) SendInitialChallenge() {

	// Initial Sending of challenge
	client.writer.Write([]byte(AUTHINFO))
	client.writer.Flush()
	client.Status = CTMS_WAITING
}

func (client *Client) DenyAccess() {
	fmt.Println("DENYING ACCESS TO CLIENT")
	client.writer.Write([]byte(AUTHDENIED)) // DENIED
	client.Status = CTMS_TERMINATED
}

func (client *Client) AllowAccess() {
	fmt.Println("APPROVING CLIENT")
	client.writer.Write([]byte(AUTH_APPROVED)) // APPROVED
	client.Status = CTMS_LOGGED
}

func (client *Client) SendMOTD() {

	fmt.Println("SENDING MOTD")
	msg := "HELLO EVERYONE HOW ARE YOU DOINGAAA\nAAAA\nAAAAAAAAA\naasfsdf\nsf\r\n"
	message := []byte(msg)
	msglen := len(msg)
	msgreq := []byte{00, 00, 00}

	MOTD_Answer := []byte{00}

	for i := 0; i < 3; i++ {
		MOTD_Answer = append(MOTD_Answer, byte(00))
	}

	MOTD_Answer = append(MOTD_Answer, byte(msglen))

	MOTD_Answer = append(MOTD_Answer, message...)

	msglen = len(MOTD_Answer)
	MOTD_Answer[0] = byte(msglen)
	MOTD_Answer = append(MOTD_Answer, msgreq...)
	MOTD_Answer = append(MOTD_Answer, byte(0))
	fmt.Println(hex.Dump(MOTD_Answer))

	client.writer.Write(MOTD_Answer) // MOTD
}

func (client *Client) ReadMOTDChallenge(buffer []byte) {

	// Only requestable by Authenticated people
	if client.Status != CTMS_LOGGED {
		return
	}

	fmt.Println("REQUESTING MOTD")

	if buffer[0] == 1 && buffer[1] == 0 && buffer[2] == 0 && buffer[3] == 0 && buffer[4] == 1 {
		client.AnswerType = MSTC_MOTD
	}
}

func (client *Client) ReadChallenge(buffer []byte, buffersize int) {

	if client.Status == CTMS_LOGGED {
		return
	}

	// Making our client verification checks
	if client.Status == CTMS_WAITING {

		if buffersize != 96 && buffersize != 100 {
			// Authentication seems unusual, drop
			client.DenyAccess()
			return
		}

		challengehash := buffer[0:71]

		// Determine the kind of server used
		switch challengehash[0] {
		case 92: // '<'
			client.ClientType = CTYPE_CLIENT
			break
		case 96: // '`'
			client.ClientType = CTYPE_SERVER
			break
		default:
			// Not a client, not a server, GET OUT.
			client.DenyAccess()
			return
		}

		// Everything is OK, approve!
		client.AnswerType = MSTC_APPROVED
	}
}

func (client *Client) Authentication_WriteData() {

	switch client.AnswerType {

	case MSTC_APPROVED:
		client.AllowAccess()
		break

	case MSTC_MOTD:
		client.SendMOTD()
	}
	client.writer.Flush()
}
