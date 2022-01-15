package main

import (
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
