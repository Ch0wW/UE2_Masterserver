package main

import (
	"encoding/hex"
	"errors"
	"fmt"
)

type ServerInfo struct {
	IP                []string
	Gameport          int32
	ServerName        string
	GameMode          string
	GameMap           string
	CurrentPlayers    byte
	MaxPlayers        byte
	IsDedicatedServer byte
	IsPrivateServer   byte

	Rules   []ServerRules
	Players []ServerPlayers
}

type ServerRules struct {
	Name  string
	Value string
}

type ServerPlayers struct {
	ID    int32
	Name  string
	Ping  int32
	Score byte
}

func (cl *UnrealConnection) Parse_PariahGameInfo() (*ServerInfo, error) {

	sv := &ServerInfo{}

	test, err := cl.ReadByte()
	if err != nil {
		return &ServerInfo{}, err
	}
	if test != 1 {
		return &ServerInfo{}, errors.New("received value should be 1")
	}

	numclients, err := cl.ReadByte()
	if err != nil {
		return &ServerInfo{}, errors.New("unable to understand players")
	}

	var i byte
	for i = 0; i < numclients-1; i++ {
		str, err := cl.ReadString()
		if err != nil {
			return &ServerInfo{}, err
		}
		sv.IP = append(sv.IP, str)
	}

	for i := 0; i < 13; i++ {
		_, err := cl.ReadByte()
		if err != nil {
			return &ServerInfo{}, err
		}
	}

	sv.Gameport, err = cl.ReadLong()
	if err != nil {
		return &ServerInfo{}, err
	}

	sv.ServerName, err = cl.ReadString()
	if err != nil {
		return &ServerInfo{}, err
	}

	sv.GameMap, err = cl.ReadString()
	if err != nil {
		return &ServerInfo{}, err
	}

	sv.GameMode, err = cl.ReadString()
	if err != nil {
		return &ServerInfo{}, err
	}

	sv.CurrentPlayers, err = cl.ReadByte()
	if err != nil {
		return &ServerInfo{}, err
	}

	sv.MaxPlayers, err = cl.ReadByte()
	if err != nil {
		return &ServerInfo{}, err
	}

	sv.IsDedicatedServer, err = cl.ReadByte()
	if err != nil {
		return &ServerInfo{}, err
	}

	sv.IsPrivateServer, err = cl.ReadByte()
	if err != nil {
		return &ServerInfo{}, err
	}

	numRules, err := cl.ReadByte()
	if err != nil {
		return &ServerInfo{}, err
	}

	for i := 0; i < int(numRules-1); i++ {

		name, err := cl.ReadString()
		if err != nil {
			return &ServerInfo{}, err
		}

		value, err := cl.ReadString()
		if err != nil {
			return &ServerInfo{}, err
		}

		a := ServerRules{
			Name:  name,
			Value: value,
		}

		sv.Rules = append(sv.Rules, a)
	}

	numPlayers, err := cl.ReadByte()
	if err != nil {
		return &ServerInfo{}, err
	}

	for i := 0; i < int(numPlayers-1); i++ {

		id, err := cl.ReadLong()
		if err != nil {
			return &ServerInfo{}, err
		}

		name, err := cl.ReadString()
		if err != nil {
			return &ServerInfo{}, err
		}

		ping, err := cl.ReadLong()
		if err != nil {
			return &ServerInfo{}, err
		}

		score, err := cl.ReadByte()
		if err != nil {
			return &ServerInfo{}, err
		}

		player := ServerPlayers{
			ID:    id,
			Name:  name,
			Ping:  ping,
			Score: score,
		}

		sv.Players = append(sv.Players, player)
	}

	fmt.Println(sv)

	return sv, nil
}

func (cl *UnrealConnection) Parse_PariahServerList() error {

	check, err := cl.ReadByte()
	if err != nil {
		return err
	}
	if check != 0 {
		return fmt.Errorf("unknown masterserver challenge code (got %d)", check)
	}

	// Check length
	packetlen := cl.bufferlen - cl.bufferpos

	// Packet is empty, give all servers
	if packetlen == 1 {
		// return a test package
	}

	/*	arguments, err := cl.ReadByte()
		if err != nil {
			return err
		}

		for */

	var pkt UnrealPacket

	pkt.WriteInt(1) // amount of servers

	// Example of how 1 server is handled
	pkt.WriteByte(1)                   // ServerID (MUST BE > 0 !!!!)
	pkt.WriteByte(127)                 // IP A
	pkt.WriteByte(0)                   // IP B
	pkt.WriteByte(0)                   // IP C
	pkt.WriteByte(1)                   // IP D
	pkt.WriteShort(7777)               // 7897 // GAME PORT
	pkt.WriteShort(7778)               // 7898 // QUERY PORT
	pkt.WriteString("HELLO WORLD FFS") // SERVERNAME
	pkt.WriteString("DM-Grind")        // MAPNAME
	pkt.WriteString("xDeathMatch")     // Gamemode
	pkt.WriteByte(0)                   // Players
	pkt.WriteByte(8)                   // Maxplayers

	fmt.Println(hex.Dump(pkt.ExportToBytes()))
	_, err = cl.conn.Write(pkt.ExportToBytes())

	if err != nil {
		return fmt.Errorf("can't send servers to client %s", cl.conn.RemoteAddr().String())
	}

	return nil
}
