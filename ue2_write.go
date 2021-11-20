package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
)

func (cl *UnrealConnection) SendInitialChallenge() error {

	cmd := "0"
	cmdlen := cl.GetStringLength(cmd)

	var buf bytes.Buffer
	buf = cl.WriteCommandSize(buf, cmdlen) // Size of command buffer
	buf = cl.WriteString(buf, cmd)

	fmt.Println(hex.Dump(buf.Bytes()))
	_, err := cl.conn.Write(buf.Bytes())

	if err != nil {
		return errors.New("Can't send initial challenge.")
	}

	//cl.Status = CTMS_WAITING

	return nil
}

func (cl *UnrealConnection) DenyAccess() error {
	fmt.Println("DENYING ACCESS TO", cl.conn.RemoteAddr().String())

	cmd := "DENIED"
	cmdlen := cl.GetStringLength(cmd)

	var buf bytes.Buffer
	buf = cl.WriteCommandSize(buf, cmdlen) // Size of command buffer
	buf = cl.WriteString(buf, cmd)
	_, err := cl.conn.Write(buf.Bytes())

	if err != nil {
		return errors.New("Cannot send it.")
	}
	return nil
}

func (cl *UnrealConnection) AllowAccess() error {
	fmt.Println("APPROVING CLIENT", cl.conn.RemoteAddr().String())

	cmd := "APPROVED"
	cmdlen := uint32(len(cmd)) + 2

	var buf bytes.Buffer
	buf = cl.WriteCommandSize(buf, cmdlen) // Size of command buffer
	buf = cl.WriteString(buf, cmd)

	fmt.Println(hex.Dump(buf.Bytes()))
	_, err := cl.conn.Write(buf.Bytes())

	if err != nil {
		return errors.New("Cannot send approving access.")
	}

	cl.Status = CTMS_LOGGED

	// Since the next logical state of Unreal Engine 2.X seems to be the UDP queries...
	if cl.Protocol.clienttype == CTYPE_SERVER {
		cl.Status = CTMS_UDPAUTHREQUEST
	}

	return nil
}

func (cl *UnrealConnection) SendVerifiedMSG() error {
	fmt.Println("Verifying CLIENT", cl.conn.RemoteAddr().String())

	cmd := "VERIFIED"
	cmdlen := uint32(len(cmd)) + 2

	var buf bytes.Buffer
	buf = cl.WriteCommandSize(buf, cmdlen) // Size of command buffer
	buf = cl.WriteString(buf, cmd)

	fmt.Println(hex.Dump(buf.Bytes()))
	_, err := cl.conn.Write(buf.Bytes())

	if err != nil {
		return errors.New("Cannot send approving access.")
	}

	cl.Status = CTMS_LOGGED

	return nil
}

func (cl *UnrealConnection) SendMOTD() error {

	var premsg, buf bytes.Buffer
	fmt.Println("SENDING MOTD")

	// MSGLINE ()
	msgline1 := "This is an Unreal Engine 2.X PoC"

	// Forge the MOTD
	premsg = cl.WriteStringNoText(premsg, msgline1)
	/*premsg.Write([]byte{0x0a, 0x0d})
	premsg = cl.WriteStringNoText(premsg, msgline2)*/
	premsg.Write([]byte{00, 00, 00, 00, 00})

	// Get the full message size for later...
	lenmsg := uint32(len(premsg.Bytes()))

	// Okay, forge the packet.
	buf = cl.WriteCommandSize(buf, lenmsg)
	buf.Write(premsg.Bytes())
	fmt.Println(hex.Dump(buf.Bytes()))

	_, err := cl.conn.Write(buf.Bytes())

	if err != nil {
		return errors.New("Cannot send it.")
	}

	return nil
}

/*
00000000  44 00 00 00 3f 48 45 4c  4c 4f 20 45 56 45 52 59  |D...?HELLO EVERY|
00000010  4f 4e 45 20 48 4f 57 20  41 52 45 20 59 4f 55 20  |ONE HOW ARE YOU |
00000020  44 4f 49 4e 47 41 41 41  0a 41 41 41 41 0a 41 41  |DOINGAAA.AAAA.AA|
00000030  41 41 41 41 41 41 41 0a  61 61 73 66 73 64 66 0a  |AAAAAAA.aasfsdf.|
*/

/*
00000000  0a 00 00 00 05 48 45 4c  4c 4f 00 00 00 00        |.....HELLO....|

*/

/*
NOTOK
00000000  0d 00 00 00 05 48 45 4c  4c 4f 00 00 00 00        |.....HELLO....|
*/
