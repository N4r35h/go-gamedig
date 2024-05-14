package gamequery

import (
	"strconv"
	"strings"

	"github.com/N4r35h/go-gamedig/types"
	ping "github.com/alteamc/minequery/ping"
)

func Minecraftquery(hostname string) (string, error) {
	gamedata := MinecraftqueryAsStruct(hostname)
	djson := jsonify(gamedata)
	return (string(djson)), nil
}

func MinecraftqueryAsStruct(hostname string) types.Gamequery {
	var gamedata types.Gamequery
	splits := strings.Split(hostname, ":")
	port, _ := strconv.Atoi(splits[1])
	res, err := ping.Ping(splits[0], uint16(port))
	if err != nil {
		gamedata.Error = "Error"
		return gamedata
	}

	gamedata.Onlineplayers = res.Players.Online
	gamedata.Maxplayers = res.Players.Max
	for _, v := range res.Players.Sample {
		gamedata.Detailed = append(gamedata.Detailed, map[string]string{
			"name": v.Name,
			"id":   v.ID,
		})
	}

	legres, err := ping.PingLegacy(splits[0], uint16(port))
	if err == nil {
		gamedata.ServerInfo.Name = legres.MessageOfTheDay
	}
	gamedata.ServerInfo.OnlinePlayers = res.Players.Online
	gamedata.ServerInfo.MaxPlayers = res.Players.Max

	gamedata.ServerInfo.Extra = map[string]string{
		"Version Name":     res.Version.Name,
		"Version Protocol": strconv.Itoa(res.Version.Protocol),
		"Favicon":          res.Favicon,
	}
	return gamedata
}
