package gamequery

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/N4r35h/go-gamedig/types"
	"github.com/sirupsen/logrus"
)

type Mtaasequery struct {
	Protocol    string
	Game        string
	Gameport    string
	Hostname    string
	Gametype    string
	Maptype     string
	Version     string
	Passworded  string
	Playercount string
	Maxplayers  string
	Rules       []Mtarule
	Players     []Mtaplayer
}

type Mtarule struct {
	Key   string
	Value string
}

type Mtaplayer struct {
	Nick  string
	Team  string
	Skin  string
	Score string
	Ping  string
	Time  string
}

func Mtaquery(hostname string) (string, error) {
	gamedata := MtaqueryAsStruct(hostname)
	djson := jsonify(gamedata)
	return (string(djson)), nil
}

func MtaqueryAsStruct(hostname string) types.Gamequery {
	host := strings.Split(hostname, ":")
	aseport, _ := strconv.Atoi(host[1])
	aseport += 123
	hostname = host[0] + ":" + strconv.Itoa(aseport)
	var gamedata types.Gamequery
	var mtaaseresp Mtaasequery
	s, _ := net.ResolveUDPAddr("udp4", hostname)
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		gamedata.Error = "Error"
		return gamedata
	}

	defer c.Close()

	data := []byte("s" + "\n")
	_, err = c.Write(data)

	if err != nil {
		logrus.Warn(err)
		gamedata.Error = "Error"
		return gamedata
	}

	buffer := make([]byte, 1024)
	n := 0

	c1 := make(chan int, 1)
	go func() {
		n, _, err = c.ReadFromUDP(buffer)
		c1 <- n
	}()

	select {
	case <-c1:
	case <-time.After(2 * time.Second):
		gamedata.Error = "Error"
		return gamedata
	}

	if err != nil {
		logrus.Warn(err)
		gamedata.Error = "Error"
		return gamedata
	}

	mtaaseresp.Protocol = string(buffer[0:4])
	mtaaseresp.Game = string(buffer[5:8])

	if mtaaseresp.Protocol != "EYE1" || mtaaseresp.Game != "mta" {
		gamedata.Error = "Error"
		return gamedata
	}
	ptr := 8

	gameportlen := int(buffer[ptr])
	ptr++
	mtaaseresp.Gameport = string(buffer[ptr : ptr+gameportlen-1])
	ptr += gameportlen - 1

	hostnamelength := int(buffer[ptr])
	ptr++
	mtaaseresp.Hostname = string(buffer[ptr : ptr+hostnamelength-1])
	ptr += hostnamelength - 1

	gametypelength := int(buffer[ptr])
	ptr++
	mtaaseresp.Gametype = string(buffer[ptr : ptr+gametypelength-1])
	ptr += gametypelength - 1

	maptypelength := int(buffer[ptr])
	ptr++
	mtaaseresp.Maptype = string(buffer[ptr : ptr+maptypelength-1])
	ptr += maptypelength - 1

	versionlength := int(buffer[ptr])
	ptr++
	mtaaseresp.Version = string(buffer[ptr : ptr+versionlength-1])
	ptr += versionlength - 1

	passwordedlength := int(buffer[ptr])
	ptr++
	mtaaseresp.Passworded = string(buffer[ptr : ptr+passwordedlength-1])
	ptr += passwordedlength - 1

	playercountlength := int(buffer[ptr])
	ptr++
	mtaaseresp.Playercount = string(buffer[ptr : ptr+playercountlength-1])
	ptr += playercountlength - 1

	maxplayerslength := int(buffer[ptr])
	ptr++
	mtaaseresp.Maxplayers = string(buffer[ptr : ptr+maxplayerslength-1])
	ptr += maxplayerslength - 1

	for {
		if int(buffer[ptr]) == 1 {
			ptr++
			break
		} else {
			rulekeylen := int(buffer[ptr])
			ptr++
			ptr += rulekeylen - 1
			rulevallen := int(buffer[ptr])
			ptr++
			ptr += rulevallen - 1
		}
	}

	pCount, _ := strconv.Atoi(mtaaseresp.Playercount)
	for i := 0; i < pCount; i++ {

		var mtaplayer Mtaplayer
		ptr++
		nicklength := int(buffer[ptr])
		ptr++
		if nicklength == 0 {
			mtaplayer.Nick = ""
		} else {
			mtaplayer.Nick = string(buffer[ptr : ptr+nicklength-1])
		}
		ptr += nicklength - 1
		teamlength := int(buffer[ptr])
		ptr++
		if teamlength == 0 {
			mtaplayer.Team = ""
		} else {
			mtaplayer.Team = string(buffer[ptr : ptr+teamlength-1])
		}
		ptr += teamlength - 1
		skinlength := int(buffer[ptr])
		ptr++
		if skinlength == 0 {
			mtaplayer.Skin = ""
		} else {
			mtaplayer.Skin = string(buffer[ptr : ptr+skinlength-1])
		}
		ptr += skinlength - 1
		scorelength := int(buffer[ptr])
		ptr++
		if scorelength == 0 {
			mtaplayer.Score = ""
		} else {
			mtaplayer.Score = string(buffer[ptr : ptr+scorelength-1])
		}
		ptr += scorelength - 1
		pinglength := int(buffer[ptr])
		ptr++
		if pinglength == 0 {
			mtaplayer.Ping = ""
		} else {
			mtaplayer.Ping = string(buffer[ptr : ptr+pinglength-1])
		}
		ptr += pinglength - 1
		timelength := int(buffer[ptr])
		ptr++
		if timelength == 0 {
			mtaplayer.Time = ""
		} else {
			mtaplayer.Time = string(buffer[ptr : ptr+timelength-1])
		}
		ptr += timelength - 1
		mtaaseresp.Players = append(mtaaseresp.Players, mtaplayer)
	}
	gamedata.Onlineplayers, _ = strconv.Atoi(mtaaseresp.Playercount)
	gamedata.Maxplayers, _ = strconv.Atoi(mtaaseresp.Maxplayers)
	for _, s := range mtaaseresp.Players {
		gamedata.Detailed = append(gamedata.Detailed, map[string]string{
			"name": s.Nick,
			"ping": s.Ping,
		})
	}

	gamedata.ServerInfo.Name = mtaaseresp.Hostname
	gamedata.ServerInfo.OnlinePlayers, _ = strconv.Atoi(mtaaseresp.Playercount)
	gamedata.ServerInfo.MaxPlayers, _ = strconv.Atoi(mtaaseresp.Maxplayers)

	gamedata.ServerInfo.Extra = map[string]string{
		"Gametype": mtaaseresp.Gametype,
		"Version":  mtaaseresp.Version,
		"Maptype":  mtaaseresp.Maptype,
		"Password": mtaaseresp.Passworded,
	}

	for _, v := range mtaaseresp.Rules {
		gamedata.ServerInfo.Extra[v.Key] = v.Value
	}

	gamedata.Error = ""
	return gamedata
}
