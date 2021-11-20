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

	if sv.bufferpos+1 > sv.bufferlen {
		errmsg := fmt.Sprintf("Buffer going too far! (pos: %d, size:%d)", sv.bufferpos+1, sv.bufferlen)
		return 0, errors.New(errmsg)
	}

	val := sv.buffer[sv.bufferpos]
	sv.bufferpos = sv.bufferpos + 1

	return val, nil
}

// ReadShort - Reads a short into the request list.
// Moves 2 bytes in the request position.
func (sv *UnrealConnection) ReadShort() (int16, error) {

	if sv.bufferpos+2 > sv.bufferlen {
		errmsg := fmt.Sprintf("Buffer going too far! (pos: %d, size:%d)", sv.bufferpos+2, sv.bufferlen)
		return 0, errors.New(errmsg)
	}

	test := binary.LittleEndian.Uint16(sv.buffer[sv.bufferpos:])
	value := uint16(test)
	sv.bufferpos = sv.bufferpos + 2

	return int16(value), nil
}

func (sv *UnrealConnection) ReadUShort() uint16 {
	test := binary.LittleEndian.Uint16(sv.buffer[sv.bufferpos:])
	value := uint16(test)
	sv.bufferpos = sv.bufferpos + 2
	return value
}

// Transform the byte into a long.
func (sv *UnrealConnection) ReadULong() (uint32, error) {

	if sv.bufferpos+4 > sv.bufferlen {
		return 0, errors.New("Buffer going too far!")
	}

	test := binary.LittleEndian.Uint32(sv.buffer[sv.bufferpos:])
	value := uint32(test)
	sv.bufferpos = sv.bufferpos + 4
	return value, nil
}

// Read long (int) bytes.
func (sv *UnrealConnection) ReadLong() (int32, error) {

	if sv.bufferpos+4 > sv.bufferlen {
		errmsg := fmt.Sprintf("Buffer going too far! (pos: %d, size:%d)", sv.bufferpos+4, sv.bufferlen)
		return 0, errors.New(errmsg)
	}

	test := binary.LittleEndian.Uint32(sv.buffer[sv.bufferpos:])
	value := uint32(test)
	sv.bufferpos = sv.bufferpos + 4

	return int32(value), nil
}

// Reads a String
/*
	To see a string, the Unreal Engine 2.X does the following:
	- Check the size of the string
	- Then read the size of the byte for the string
	- It ends with "0x00"
*/
func (sv *UnrealConnection) ReadString() (string, error) {

	result := ""
	strsize, err := sv.ReadByte()
	if err != nil {
		return "", err
	}

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

	fmt.Println("Size of String:", strsize-1, "==>", result)
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

func (msg *UnrealConnection) WriteCommandSize(buf bytes.Buffer, packetsize uint32) bytes.Buffer {

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, packetsize)
	buf.Write(b)
	return buf
}

func (msg *UnrealConnection) GetStringLength(cmd string) uint32 {
	return uint32(len(cmd) + 2) // Size + 00
}
