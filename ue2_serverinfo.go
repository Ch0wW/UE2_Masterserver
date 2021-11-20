package main

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
)

func (msg *UnrealConnection) Server_GetUDPPortRequest() error {

	_, err := msg.ReadLong()
	if err != nil {
		return errors.New("Unknown value One")
	}

	_, err = msg.ReadLong()
	if err != nil {
		fmt.Println(err)
		return errors.New("Unknown value Two")
	}

	var buf bytes.Buffer
	msg.SV_UDPInfo.Code = uint32(rand.Intn((10000 - 1) + 1)) // Challenge key...

	// Send UDP queries data...
	buf.Write([]byte{06, 00, 00, 00})
	buf.WriteByte(0)
	buf.WriteByte(0)
	buf = msg.WriteCommandSize(buf, msg.SV_UDPInfo.Code)

	buf.Write([]byte{06, 00, 00, 00})
	buf.WriteByte(0)
	buf.WriteByte(1)
	buf = msg.WriteCommandSize(buf, msg.SV_UDPInfo.Code)

	// Basically Gamespy died but... Yeah, "in case of"?
	buf.Write([]byte{06, 00, 00, 00})
	buf.WriteByte(0)
	buf.WriteByte(2)
	buf = msg.WriteCommandSize(buf, msg.SV_UDPInfo.Code)

	_, err = msg.conn.Write(buf.Bytes())

	if err != nil {
		return errors.New("Cannot send it.")
	}

	msg.Status = CTMS_SERVAUTH02
	return nil
}

func (msg *UnrealConnection) Server_GetServerInfoRequest() error {

	bytecmd, err := msg.ReadByte()
	if err != nil {
		return errors.New("Unknown value One")
	}

	switch bytecmd {
	case 3:

		var buf bytes.Buffer
		buf.Write([]byte{04, 00, 00, 00})
		buf.Write([]byte{04, 00, 00, 00})

		fmt.Println("Sending", buf.Bytes())

		_, err = msg.conn.Write(buf.Bytes())

		if err != nil {
			return errors.New("Cannot send it.")
		}
	}

	return nil
}
