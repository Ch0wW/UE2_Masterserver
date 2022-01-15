package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

func (cl *UnrealConnection) SendSimpleString(text string) error {
	var pkt UnrealPacket
	pkt.WriteString(text)

	if debugmode {
		fmt.Println("Sending", text, "string packet to client", cl.conn.RemoteAddr().String())
		fmt.Println(hex.Dump(pkt.ExportToBytes()))
	}

	fmt.Println(hex.Dump(pkt.ExportToBytes()))
	_, err := cl.conn.Write(pkt.ExportToBytes())

	if err != nil {
		return fmt.Errorf("can't send message %s to client %s", text, cl.conn.RemoteAddr().String())
	}

	return nil
}

func (cl *UnrealConnection) SendVerifiedMSG() error {

	err := cl.SendSimpleString("VERIFIED")
	if err != nil {
		return err
	}

	// Should validate the client
	cl.Status = CTMS_LOGGED

	return nil
}

func (cl *UnrealConnection) SendMOTD() error {

	var premsg, buf bytes.Buffer
	fmt.Println("SENDING MOTD TO CLIENT")

	// MSGLINE ()
	msgline1 := "Bonjour Epic Games! Ici Ch0wW."

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
		return errors.New("unable to send MOTD message")
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

type UnrealPacket struct {
	buf bytes.Buffer // Buffer to send
}

func (pkt *UnrealPacket) WriteInt(packetsize uint32) {

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, packetsize)
	pkt.buf.Write(b)
}

func (pkt *UnrealPacket) WriteString(cmd string) {
	pkt.buf.WriteByte(byte(len(cmd)) + 1)
	pkt.buf.Write([]byte(cmd))
	pkt.buf.WriteByte(0)
}

func (pkt *UnrealPacket) WriteByte(c byte) error {

	err := pkt.buf.WriteByte(c)
	if err != nil {
		return err
	}
	return nil
}

func (pkt *UnrealPacket) Write(c []byte) (int, error) {

	n, err := pkt.buf.Write(c)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (pkt *UnrealPacket) ExportToBytes() []byte {

	var retval UnrealPacket

	size := len(pkt.buf.Bytes())
	retval.WriteInt(uint32(size))

	retval.buf.Write(pkt.buf.Bytes())

	return retval.buf.Bytes()
}

func (pkt *UnrealPacket) Reset() {
	pkt.buf.Reset()
}
