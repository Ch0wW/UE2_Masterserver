package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var (
	botConfig Config
)

type LatestVersion struct {
	/*DOOM3  LatestVersionDetails `json:"doom3"`
	DHEWM3 LatestVersionDetails `json:"dhewm3"`
	QUAKE4 LatestVersionDetails `json:"quake4"`
	PREY   LatestVersionDetails `json:"prey"`*/
}

type Config struct {
	Games LatestVersion `json:"gameinfo"`
}

func BotConfig_Init() {
	fmt.Println("> Loading and parsing config...")

	cfg, err := os.Open("config.json")
	if err != nil {
		log.Printf("[error] failed to open config.json: %v\n", err)
		os.Exit(1)
	}

	err = json.NewDecoder(cfg).Decode(&botConfig)
	if err != nil {
		log.Printf("[error] failed to decode config.json: %v\n", err)
		os.Exit(1)
	}
}

/*
func (cl *UnrealConnection) SendMOTD() error {

	fmt.Println("SENDING MOTD")
	//	msg := "HELLO EVERYONE HOW ARE YOU DOINGAAA\nAAAA\nAAAAAAAAA\naasfsdf\nsf\n"
	msg := "HELLO"

	var premsg, buf bytes.Buffer
	premsg = cl.WritePayload(premsg)
	premsg.WriteByte(byte(len(msg)))
	premsg.Write([]byte(msg))
	premsg.Write([]byte{00, 00, 00, 00})

	buf.WriteByte(byte(len(premsg.Bytes())))
	buf.Write(premsg.Bytes())
	fmt.Println(hex.Dump(buf.Bytes()))

	_, err := cl.conn.Write(buf.Bytes())

	if err != nil {
		return errors.New("Cannot send it.")
	}

	return nil
}*/
