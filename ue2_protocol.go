package main

type Protocol int

const (
	PROTOCOL_NONE int = iota
	PROTOCOL_PARIAH
	PROTOCOL_LOTD
	PROTOCOL_UT2KX
)

type ClientType int

const (
	CTYPE_UNKNOWN ClientType = iota
	CTYPE_CLIENT
	CTYPE_SERVER
)

type UE2_GameProtocol struct {
	protocol   Protocol
	clienttype ClientType
}

func PROTOCOL_Init() {

}