package gamequery

import (
	"encoding/binary"
	"net"
	"strconv"
	"time"

	"github.com/N4r35h/go-gamedig/types"
)

type VCMPInfo struct {
	Protocol    string
	Game        string
	Gameport    string
	Hostname    string
	Gametype    string
	Maptype     string
	Version     string
	Passworded  string
	Playercount int
	Maxplayers  int
	Players     []string
}

func VCMPqueryAsStruct(hostname string) types.Gamequery {
	var gamedata types.Gamequery
	var vcmpresponse VCMPInfo

	buffer, err := VCMPCoreQuery(hostname, "i")
	if err != nil {
		gamedata.Error = "Core fail"
		return gamedata
	}

	if len(buffer) == 0 {
		gamedata.Error = "Server is offline."
		return gamedata
	}
	vcmpresponse.Protocol = string(buffer[0:4])
	vcmpresponse.Game = "VCMP" + string(buffer[10])

	ptr := 11
	vcmpresponse.Version = string(buffer[ptr : ptr+11])
	ptr += 11
	vcmpresponse.Passworded = strconv.FormatBool(buffer[ptr] == 1)
	vcmpresponse.Playercount = int(binary.LittleEndian.Uint16(buffer[ptr : ptr+4]))
	ptr += 4
	vcmpresponse.Maxplayers = int(binary.LittleEndian.Uint16(buffer[ptr : ptr+2]))
	ptr += 2
	lenservername := int(binary.LittleEndian.Uint16(buffer[ptr : ptr+4]))
	ptr += 4
	vcmpresponse.Hostname = string(buffer[ptr : ptr+lenservername])
	ptr += lenservername
	lengamemodename := int(binary.LittleEndian.Uint16(buffer[ptr : ptr+4]))
	ptr += 4
	vcmpresponse.Gametype = string(buffer[ptr : ptr+lengamemodename])
	ptr += lengamemodename
	lenmapname := int(binary.LittleEndian.Uint16(buffer[ptr : ptr+4]))
	ptr += 4
	vcmpresponse.Maptype = string(buffer[ptr : ptr+lenmapname])
	ptr += lenmapname

	buffer, err = VCMPCoreQuery(hostname, "c")
	if err != nil {
		gamedata.Error = "c query fail"
		return gamedata
	}
	noofplayers := int(binary.LittleEndian.Uint16(buffer[11:13]))
	vcmpresponse.Playercount = noofplayers
	ptr = 13
	for i := 1; i <= noofplayers; i++ {
		lenplayername := int(buffer[ptr])
		ptr += 1
		playername := string(buffer[ptr : ptr+lenplayername])
		vcmpresponse.Players = append(vcmpresponse.Players, playername)
		ptr += lenplayername
	}
	gamedata.Maxplayers = vcmpresponse.Maxplayers
	gamedata.Onlineplayers = vcmpresponse.Playercount
	for _, s := range vcmpresponse.Players {
		gamedata.Detailed = append(gamedata.Detailed, map[string]string{
			"name": s,
		})
	}

	gamedata.ServerInfo.Name = vcmpresponse.Hostname
	gamedata.ServerInfo.OnlinePlayers = vcmpresponse.Playercount
	gamedata.ServerInfo.MaxPlayers = vcmpresponse.Maxplayers
	return gamedata
}

func VCMPCoreQuery(hostname string, opcode string) ([]byte, error) {
	s, err := net.ResolveUDPAddr("udp4", hostname)
	if err != nil {
		return []byte{}, err
	}
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		return []byte{}, err
	}

	defer c.Close()

	data := []byte("VCMP3OÑ:c" + opcode)
	_, err = c.Write(data)

	if err != nil {
		return []byte{}, err
	}

	buffer := make([]byte, 1024)
	n := 0

	done := make(chan int, 1)
	go func() {
		n, _, err = c.ReadFromUDP(buffer)
		done <- n
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		return []byte{}, err
	}

	if err != nil {
		return []byte{}, err
	}

	return buffer, nil
}
