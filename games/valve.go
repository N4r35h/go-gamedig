package gamequery

import (
	"strconv"
	"strings"

	"github.com/FlowingSPDG/go-steam"
	"github.com/N4r35h/go-gamedig/types"
	"github.com/sirupsen/logrus"
)

func ValveQuery(hostname string) (string, error) {
	gamedata := ValveQueryAsStruct(hostname)
	djson := jsonify(gamedata)
	return (string(djson)), nil
}

func ValveQueryAsStruct(hostname string) types.Gamequery {
	var gamedata types.Gamequery

	server, err := steam.Connect(hostname)
	if err != nil {
		gamedata.Error = "Error"
		return gamedata
	}
	defer server.Close()
	playersInfo, err := server.PlayersInfo()
	if err != nil {
		logrus.Warn(err)
		gamedata.Error = "Error"
		return gamedata
	}
	info, err := server.Info()
	if err != nil {
		logrus.Warn(err)
		gamedata.Error = "Error"
		return gamedata
	}
	gamedata.Error = ""
	gamedata.Maxplayers = info.MaxPlayers

	playercount := 0
	for _, p := range playersInfo.Players {
		if len(p.Name) != 1 {
			gamedata.Detailed = append(gamedata.Detailed, map[string]string{
				"name":     strings.TrimSuffix(p.Name, "\u0000"),
				"duration": FormatTime(p.Duration),
				"score":    strconv.Itoa(p.Score),
			})
			playercount++
		}
	}
	gamedata.Onlineplayers = playercount

	gamedata.ServerInfo.Name = info.Name
	gamedata.ServerInfo.OnlinePlayers = playercount
	gamedata.ServerInfo.MaxPlayers = info.MaxPlayers

	gamedata.Error = ""
	return gamedata
}

func FormatTime(sec float64) string {
	resultant := ""
	secInt := int(sec)
	if secInt/3600 >= 1 {
		resultant += strconv.Itoa(secInt/3600) + "h"
		secInt = secInt % 3600
	}
	if secInt/60 >= 1 {
		resultant += strconv.Itoa(secInt/60) + "m"
	}
	return resultant
}
