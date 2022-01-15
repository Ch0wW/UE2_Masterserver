package main

import (
	"errors"
	"fmt"
)

type ConnectionStatus int

const (
	CTMS_NEEDSAUTH ConnectionStatus = iota
	CTMS_WAITING
	CTMS_TERMINATED
	CTMS_NEEDSVERIFICATION
	CTMS_LOGGED

	SVTMS_UDPAUTHREQUEST // Servers only, to auth with the masterserver
	CTMS_SERVAUTH02
	SVTMS_PROCESSQUERYMESSAGE
)

type SV_AnswerType int

const (
	MSTC_UNKNOWN SV_AnswerType = iota
	MSTC_DENIED
	MSTC_APPROVED
	MSTC_UPGRADE // Not Set
	MSTC_MSLIST  // Not Set
	MSTC_MOTD
)

func (cl *UnrealConnection) SetGameType(val string) {

	switch val {
	case "PARIAHCLIENT":
		cl.Protocol.protocol = Protocol(PROTOCOL_PARIAH)
		cl.Protocol.clienttype = CTYPE_CLIENT
	case "PARIAHSERVER":
		cl.Protocol.protocol = Protocol(PROTOCOL_PARIAH)
		cl.Protocol.clienttype = CTYPE_SERVER
	case "SERVER":
		// ToDo: we'll have to do something as several games use this
		cl.Protocol.protocol = Protocol(PROTOCOL_GENERIC)
		cl.Protocol.clienttype = CTYPE_SERVER
	}

}

func (cl *UnrealConnection) ProcessAuthMessage() error {

	// CD-Key
	_, err := cl.ReadString()
	if err != nil {
		fmt.Println(err)
		return errors.New("unable to read CDKey")
	}

	// Unknown
	_, err = cl.ReadString()
	if err != nil {
		fmt.Println(err)
		return errors.New("unable to read Unknown Value #1")
	}

	// GameType
	Gametype, err := cl.ReadString()
	if err != nil {
		fmt.Println(err)
		return errors.New("unable to read Gametype")
	}
	cl.SetGameType(Gametype)

	// Game version...
	GameVersion, err := cl.ReadLong()
	if err != nil {
		fmt.Println(err)
		return errors.New("unable to read GameVersion")
	}
	fmt.Println("Game version is", GameVersion)

	// Platform
	_, err = cl.ReadByte()
	if err != nil {
		fmt.Println(err)
		return errors.New("unable to read Platform")
	}

	if cl.Protocol.clienttype == CTYPE_SERVER {
		// Unknown too...
		_, err = cl.ReadLong()
		if err != nil {
			fmt.Println(err)
			return errors.New("unable to read Server Unknown Value ")
		}

		// Initialize the serverinfo code...
		cl.SV_UDPInfo = &ServerUDPRequest{
			Code:             0,
			GamePort:         0,
			GameQueryPort:    0,
			GameSpyQueryPort: 0,
		}
	}

	lang, err := cl.ReadString()
	if err != nil {
		fmt.Println(err)
		return errors.New("unable to read Language value")
	}
	cl.Protocol.language = lang

	// Past this, client seems OK!
	return nil
}

func (cl *UnrealConnection) ProcessMOTDRequest() error {

	recvbyte, err := cl.ReadByte()
	if err != nil {
		return err
	}

	if recvbyte == 1 {
		err := cl.SendMOTD()
		if err != nil {
			return errors.New("unable to send MOTD message")
		}
	} else {
		return fmt.Errorf("unknown single byte %d sent by the client", recvbyte)
	}

	return nil
}

func (cl *UnrealConnection) ProcessKeyVerification() error {

	cl.SendVerifiedMSG()
	return nil
}

func (cl *UnrealConnection) ReadMessage() error {

	// Auth status
	if cl.Status == CTMS_NEEDSAUTH {

		err := cl.ProcessAuthMessage()
		if err != nil {
			err = cl.SendSimpleString("DENIED")
			if err != nil {
				return errors.New("this is truly problematic if the server can't even send this DENIED packet")
			}
			return errors.New("client message looks invalid")
		}

		err = cl.SendSimpleString("APPROVED")
		if err != nil {
			return errors.New("unable to send APPROVED message")
		}

		cl.Status = CTMS_LOGGED

		// ToDo: Check if verified flag
		// If we're having a game that requires the validation code... Let's do it
		if cl.Protocol.protocol == Protocol(PROTOCOL_UT2KX) {
			cl.Status = CTMS_NEEDSVERIFICATION
		} else {
			cl.Status = CTMS_LOGGED
		}

		// Since the next logical state of Unreal Engine 2.X seems to be the UDP queries...
		if cl.Protocol.clienttype == CTYPE_SERVER {
			cl.Status = SVTMS_UDPAUTHREQUEST
		}

		return nil
	}

	buflen := len(cl.buffer)

	if cl.Protocol.clienttype == CTYPE_SERVER {

		if cl.Status == SVTMS_UDPAUTHREQUEST {
			return cl.Server_GetUDPPortRequest()
		} else if cl.Status == CTMS_SERVAUTH02 {
			return cl.Server_GetServerInfoRequest()
		} else if cl.Status == SVTMS_PROCESSQUERYMESSAGE {
			cl.Parse_PariahGameInfo()
		}

		return nil
	}

	// CLIENT ONLY
	if buflen == 1 {
		return cl.ProcessMOTDRequest() // Checking if there's a MOTD
	} else if buflen == 34 && cl.Status == CTMS_NEEDSVERIFICATION {
		return cl.ProcessKeyVerification() // Checking if it's another (unknown) hash
	}

	return nil
}
