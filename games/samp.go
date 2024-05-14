package gamequery

import (
	"encoding/binary"
	"net"
	"strconv"
	"time"

	"github.com/N4r35h/go-gamedig/types"
	"github.com/sirupsen/logrus"
)

func Sampquery(hostname string) (string, error) {
	gamedata := SampqueryAsStruct(hostname)
	djson := jsonify(gamedata)
	return (string(djson)), nil
}

func SampqueryAsStruct(hostname string) types.Gamequery {
	var gamedata types.Gamequery
	gamedata.Error = ""

	resp, err := SAMPCoreQuery(hostname, "i")
	if err != nil || len(resp) == 0 {
		logrus.Warn(err)
		gamedata.Error = "unable to query server, must be offline or unreachable"
		return gamedata
	}

	ptr := 11
	passworded := false
	if int(ptr) > 0 {
		passworded = true
	}
	ptr++
	pCount := int(binary.LittleEndian.Uint16(resp[ptr : ptr+2]))
	ptr += 2
	pMaxCount := int(binary.LittleEndian.Uint16(resp[ptr : ptr+2]))
	ptr += 2
	hnLen := int(binary.LittleEndian.Uint16(resp[ptr : ptr+4]))
	ptr += 4
	servername := string(resp[ptr : ptr+hnLen])
	ptr += hnLen
	gmLen := int(binary.LittleEndian.Uint16(resp[ptr : ptr+4]))
	ptr += 4
	gamemode := string(resp[ptr : ptr+gmLen])
	ptr += gmLen
	langLen := int(binary.LittleEndian.Uint16(resp[ptr : ptr+4]))
	ptr += 4
	language := string(resp[ptr : ptr+langLen])
	ptr += langLen
	gamedata.ServerInfo.Extra = map[string]string{
		"Gamemode": gamemode,
		"Language": language,
		"Password": strconv.FormatBool(passworded),
	}
	gamedata.Onlineplayers = pCount
	gamedata.Maxplayers = pMaxCount
	gamedata.ServerInfo.Name = servername
	gamedata.ServerInfo.OnlinePlayers = pCount
	gamedata.ServerInfo.MaxPlayers = pMaxCount

	resp, err = SAMPCoreQuery(hostname, "r")
	if err != nil || len(resp) == 0 {
		logrus.Warn(err)
	} else {
		ptr := 11
		ptr += 2
	}

	resp, err = SAMPCoreQuery(hostname, "d")
	if err != nil || len(resp) == 0 {
		logrus.Warn(err)
	} else {
		ptr := 11
		pCount := int(binary.LittleEndian.Uint16(resp[ptr : ptr+2]))
		ptr += 2
		for i := 1; i <= pCount; i++ {
			pID := int(resp[ptr])
			ptr++
			pLen := int(resp[ptr])
			ptr++
			playername := string(resp[ptr : ptr+pLen])
			ptr += pLen
			score := int(binary.LittleEndian.Uint16(resp[ptr : ptr+4]))
			ptr += 4
			ping := int(binary.LittleEndian.Uint16(resp[ptr : ptr+4]))
			ptr += 4
			gamedata.Detailed = append(gamedata.Detailed, map[string]string{
				"id":    strconv.Itoa(pID),
				"name":  playername,
				"score": strconv.Itoa(score),
				"ping":  strconv.Itoa(ping),
			})
		}
	}
	return gamedata
}

func SAMPCoreQuery(hostname string, opcode string) ([]byte, error) {
	s, _ := net.ResolveUDPAddr("udp4", hostname)
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		return []byte{}, err
	}

	defer c.Close()

	data := []byte("SAMP3OÑ:c" + opcode)
	_, err = c.Write(data)

	if err != nil {
		return []byte{}, err
	}

	buffer := make([]byte, 6144)
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

	if n == 0 {
		return []byte{}, err
	}
	if err != nil {
		return []byte{}, err
	}

	return buffer, nil
}
