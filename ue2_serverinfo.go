package main

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"
)

func (msg *UnrealConnection) Server_GetUDPPortRequest() error {

	var CODE_FAILURE byte = 0
	var CODE_SUCCESS byte = 1

	var QUERY_QUERYPORT byte = 0
	var QUERY_GAMEPORT byte = 1
	var QUERY_GAMESPYPORT byte = 2

	_, err := msg.ReadLong()
	if err != nil {
		return errors.New("unknown value ONE")
	}

	gamespy, err := msg.ReadLong()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("unknown value received for gamespy: %s", err)
	}

	msg.SV_UDPInfo.Code = uint32(rand.Intn((10000 - 1) + 1)) // Challenge key...

	server := &UnrealServerData{
		IP:               msg.conn.RemoteAddr().(*net.TCPAddr).IP,
		Code:             msg.SV_UDPInfo.Code,
		Port:             0,
		QueryPort:        0,
		GameSpyQueryPort: 0,
	}

	ServerList = append(ServerList, server)

	// Send UDP queries data...
	var pkt UnrealPacket
	var buf bytes.Buffer

	pkt.WriteByte(CODE_FAILURE)
	pkt.WriteByte(QUERY_QUERYPORT)
	pkt.WriteInt(msg.SV_UDPInfo.Code)
	buf.Write(pkt.ExportToBytes())

	pkt.Reset()
	pkt.WriteByte(CODE_FAILURE)
	pkt.WriteByte(QUERY_GAMEPORT)
	pkt.WriteInt(msg.SV_UDPInfo.Code)
	buf.Write(pkt.ExportToBytes())

	if gamespy == 1 {
		pkt.Reset()
		pkt.WriteByte(CODE_FAILURE)
		pkt.WriteByte(QUERY_GAMESPYPORT)
		pkt.WriteInt(msg.SV_UDPInfo.Code)
		buf.Write(pkt.ExportToBytes())
	}

	_, err = msg.conn.Write(buf.Bytes())

	if err != nil {
		return errors.New("unable to send query ports")
	}

	// Wait to get our info...
	time.Sleep(3 * time.Second)

	for a := 0; a < len(ServerList); a++ {

		sv := ServerList[a]

		if sv.IP.Equal(msg.conn.RemoteAddr().(*net.TCPAddr).IP) {
			if sv.Code == msg.SV_UDPInfo.Code {
				var pkt UnrealPacket
				pkt.WriteByte(CODE_SUCCESS)                   // Success!
				pkt.WriteInt(uint32(botConfig.HeartBeatTime)) // Configured by the config
				pkt.WriteInt(uint32(sv.QueryPort))
				pkt.WriteInt(uint32(sv.Port))
				pkt.WriteInt(uint32(sv.GameSpyQueryPort))

				if debugmode {
					fmt.Println("Sending............", pkt.ExportToBytes())
				}

				_, err = msg.conn.Write(pkt.ExportToBytes())
				if err != nil {
					return errors.New("unable to send query status message")
				}
			}
		}

	}

	msg.Status = CTMS_SERVAUTH02
	return nil
}

func (msg *UnrealConnection) Server_GetServerInfoRequest() error {

	bytecmd, err := msg.ReadByte()
	if err != nil {
		return errors.New("unknown value ONE (ServerInfoRequest)")
	}

	switch bytecmd {
	case 4:

		var pkt UnrealPacket
		pkt.buf.WriteByte(3) // MTS_MatchID
		pkt.WriteInt(1)      // ToDo : ADD MatchID !!!
		fmt.Println("Sending", pkt.ExportToBytes())

		_, err = msg.conn.Write(pkt.ExportToBytes())

		if err != nil {
			return errors.New("cannot send it")
		}
	}

	msg.Status = SVTMS_PROCESSQUERYMESSAGE

	return nil
}
