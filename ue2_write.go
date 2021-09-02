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
	buf.WriteByte(cmdlen) // Size of Buffer
	buf = cl.WritePayload(buf)
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
	fmt.Println("DENYING ACCESS TO")
	cmd := "DENIED"
	cmdlen := cl.GetStringLength(cmd)

	var buf bytes.Buffer
	buf.WriteByte(cmdlen) // Size of Buffer
	buf = cl.WritePayload(buf)
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
	cmdlen := byte(len(cmd)) + 2

	var buf bytes.Buffer
	buf.WriteByte(cmdlen) // Size of Buffer
	buf = cl.WritePayload(buf)
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

	fmt.Println("SENDING MOTD")
	//	msg := "HELLO EVERYONE HOW ARE YOU DOINGAAA\nAAAA\nAAAAAAAAA\naasfsdf\nsf\n"
	msg := "HELLO\n"

	var premsg, buf bytes.Buffer

	//premsg.WriteByte(byte(len(msg)))
	premsg = cl.WriteStringNoText(premsg, msg)
	premsg.Write([]byte{00, 00, 00})

	buf.WriteByte(byte(len(premsg.Bytes())))
	buf = cl.WritePayload(buf)
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
