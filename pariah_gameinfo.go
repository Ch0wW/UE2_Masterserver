package main

import "errors"

type ServerInfo struct {
	IP      []string
	players byte
}

func (cl *UnrealConnection) Parse_PariahGameInfo() (*ServerInfo, error) {

	var sv *ServerInfo

	test, err := cl.ReadByte()
	if err != nil {
		return &ServerInfo{}, errors.New("unable to understand packet")
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
			return &ServerInfo{}, errors.New("Unable to understand packet")
		}
		sv.IP = append(sv.IP, str)
	}

	return sv, nil
}
