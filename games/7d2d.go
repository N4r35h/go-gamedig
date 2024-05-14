package gamequery

import (
	"strconv"
	"strings"

	"github.com/N4r35h/go-gamedig/types"

	"github.com/FlowingSPDG/go-steam"
	"github.com/sirupsen/logrus"
)

func SevenD2DQueryAsStruct(hostname string) types.Gamequery {
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
	gamedata.Maxplayers = info.MaxPlayers
	gamedata.Onlineplayers = info.Players

	playercount := 0
	for _, p := range playersInfo.Players {
		gamedata.Detailed = append(gamedata.Detailed, map[string]string{
			"name":     strings.TrimSuffix(p.Name, "\u0000"),
			"duration": FormatTime(p.Duration),
			"score":    strconv.Itoa(p.Score),
		})
		playercount++
	}

	gamedata.ServerInfo.Name = strings.TrimSuffix(info.Name, "\u0000")
	gamedata.ServerInfo.OnlinePlayers = info.Players
	gamedata.ServerInfo.MaxPlayers = info.MaxPlayers

	gamedata.ServerInfo.Extra = map[string]string{
		"Protocol":    strconv.Itoa(info.Protocol),
		"Map":         strings.Trim(info.Map, "\u0000"),
		"Folder":      strings.Trim(info.Folder, "\u0000"),
		"Game":        strings.Trim(info.Game, "\u0000"),
		"ID":          strconv.Itoa(info.ID),
		"Bots":        strconv.Itoa(info.Bots),
		"ServerType":  info.ServerType.String(),
		"Environment": info.Environment.String(),
		"Visibility":  info.Visibility.String(),
		"VAC":         strings.Trim(info.VAC.String(), "\u0000"),
		"Version":     strings.Trim(info.Version, "\u0000"),

		"Port":    strconv.Itoa(info.Port),
		"SteamID": strconv.Itoa(int(info.SteamID)),

		"SourceTVPort": strconv.Itoa(info.SourceTVPort),
		"SourceTVName": info.SourceTVName,

		"Keywords": strings.Trim(info.Keywords, "\u0000"),
		"GameID":   strconv.Itoa(int(info.GameID)),
	}

	return gamedata
}
