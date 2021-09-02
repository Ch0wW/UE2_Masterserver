package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// ReadByte - Reads the byte.
// Moves 1 byte in the request position.
func (sv *UnrealConnection) ReadByte() (byte, error) {
	//value := uint8(sv.buffer[sv.bufferpos])
	//sv.bufferpos = sv.bufferpos + 1

	if sv.bufferpos > sv.bufferlen {
		return 0, errors.New("Buffer going too far!")
	}

	val := sv.buffer[sv.bufferpos]
	sv.bufferpos = sv.bufferpos + 1

	return val, nil
}

// ReadShort - Reads a short into the request list.
// Moves 2 bytes in the request position.
func (sv *UnrealConnection) ReadShort() int16 {
	test := binary.LittleEndian.Uint16(sv.buffer[sv.bufferpos:])
	value := uint16(test)
	sv.bufferpos = sv.bufferpos + 2
	return int16(value)
}

func (sv *UnrealConnection) ReadUShort() uint16 {
	test := binary.LittleEndian.Uint16(sv.buffer[sv.bufferpos:])
	value := uint16(test)
	sv.bufferpos = sv.bufferpos + 2
	return value
}

// Transform the byte into a long.
func (sv *UnrealConnection) ReadULong() uint32 {
	test := binary.LittleEndian.Uint32(sv.buffer[sv.bufferpos:])
	value := uint32(test)
	sv.bufferpos = sv.bufferpos + 4
	return value
}

// Transform the byte into a long.
func (sv *UnrealConnection) ReadLong() int32 {
	test := binary.LittleEndian.Uint32(sv.buffer[sv.bufferpos:])
	value := uint32(test)
	sv.bufferpos = sv.bufferpos + 4
	return int32(value)
}

// Transform the byte into a long.
func (sv *UnrealConnection) ReadString() (string, error) {

	result := ""

	strsize, err := sv.ReadByte()
	if err != nil {
		return "", err
	}

	fmt.Println("Size of String:", strsize)

	for i := 0; i < int(strsize); i++ {
		c, err := sv.ReadByte()

		if err != nil {
			return "", err
		}

		if c <= 0 || c >= 255 {
			break
		}

		if c == '%' {
			c = '.'
		}

		result = result + string(c)
	}

	return result, nil
}

func (msg *UnrealConnection) WriteString(buf bytes.Buffer, cmd string) bytes.Buffer {
	buf.WriteByte(byte(len(cmd)) + 1)
	buf.Write([]byte(cmd))
	buf.WriteByte(0)

	return buf
}

func (msg *UnrealConnection) WriteStringNoText(buf bytes.Buffer, cmd string) bytes.Buffer {
	buf.WriteByte(byte(len(cmd)))
	buf.Write([]byte(cmd))
	buf.WriteByte(0)

	return buf
}

func (msg *UnrealConnection) WritePayload(buf bytes.Buffer) bytes.Buffer {
	buf.Write([]byte{00, 00, 00})
	return buf
}

func (msg *UnrealConnection) GetStringLength(cmd string) byte {
	return byte(len(cmd) + 2) // Size + 00
}
