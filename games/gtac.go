package gamequery

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/N4r35h/go-gamedig/types"
	"github.com/sirupsen/logrus"
)

func GTACqueryAsStruct(hostname string) types.Gamequery {
	response, err := GTACCoreQuery(hostname, "")
	if err != nil {
		logrus.Warn(err)
	}
	var gamedata types.Gamequery
	if len(response) == 0 {
		gamedata.Error = "Server is offline"
		return gamedata
	}
	if string(response[2:5]) != "UGP" {
		gamedata.Error = "Invalid response indentifier"
		return gamedata
	}

	SrvGameID := int(response[8])
	ptr := 11
	servernamelen := int(response[ptr]) // Server Name Lenght?
	ptr++
	SrvName := string(response[ptr : ptr+servernamelen]) // Server Name?
	ptr += servernamelen
	modelen := int(response[ptr]) // Server Mode Lenght?
	ptr++
	SrvGameMode := string(response[ptr : ptr+modelen]) // Server Mode?
	ptr += modelen
	gamedata.Onlineplayers = int(response[ptr]) // Players Count
	ptr++
	gamedata.Maxplayers = int(response[ptr]) // Maxplayers
	ptr++
	ptr += 2
	for i := 1; i <= gamedata.Onlineplayers; i++ {
		playerid := int(response[ptr]) //ID
		ptr++
		nicklen := int(response[ptr]) //Nick Len
		ptr++
		playername := string(response[ptr : ptr+nicklen-1])
		playername = strings.Replace(playername, "\u0000", "", 1)
		gamedata.Detailed = append(gamedata.Detailed, map[string]string{
			"id":   strconv.Itoa(playerid),
			"name": playername,
		})
		ptr += nicklen
	}

	gamedata.ServerInfo.Name = strings.Trim(SrvName, "\u0000")
	gamedata.ServerInfo.OnlinePlayers = gamedata.Onlineplayers
	gamedata.ServerInfo.MaxPlayers = gamedata.Maxplayers

	gamedata.ServerInfo.Extra = map[string]string{
		"Gamemode": strings.Trim(SrvGameMode, "\u0000"),
		"GameID":   strconv.Itoa(SrvGameID),
	}

	return gamedata
}

func GTACCoreQuery(hostname string, opcode string) ([]byte, error) {
	s, err := net.ResolveUDPAddr("udp4", hostname)
	if err != nil {
		return []byte{}, err
	}
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		return []byte{}, err
	}

	defer c.Close()
	data := []byte{0xff, 0xff, 0x55, 0x47, 0x50, 0x00, 0x01, 0x01, 0x03, 0x6c, 0x7f, 0x3d}
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
