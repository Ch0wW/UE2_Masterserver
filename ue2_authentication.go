package main

import (
	"fmt"
)

type ConnectionStatus int

const (
	CTMS_NEEDSAUTH ConnectionStatus = iota
	CTMS_WAITING
	CTMS_TERMINATED
	CTMS_LOGGED
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

func (cl *UnrealConnection) ReadMessage() {

	// Auth status
	if cl.Status == CTMS_NEEDSAUTH {
		_, err := cl.ReadString()

		err = cl.AllowAccess()

		if err != nil {
			cl.CloseConnection()
		}

		return
	}

	// Checking if there's a MOTD
	if len(cl.buffer) == 1 {
		fmt.Println("AAA")
		motdauth, err := cl.ReadByte()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("BYTE SAYS ==>", motdauth)

		if motdauth == 1 {
			err := cl.SendMOTD()
			if err != nil {
				fmt.Println("AAA")
			}
		}
	}
}
