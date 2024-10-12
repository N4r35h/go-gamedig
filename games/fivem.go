package gamequery

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/N4r35h/go-gamedig/types"
	"github.com/sirupsen/logrus"
)

type fiveminfojson struct {
	Resources []string        `json:"resources,omitempty"`
	Server    string          `json:"server,omitempty"`
	Icon      string          `json:"icon,omitempty"`
	Vars      fivemservervars `json:"vars,omitempty"`
	Version   int             `json:"version,omitempty"`
}

type fivemservervars struct {
	Banner_detail   string `json:"banner_detail,omitempty"`
	GameBuild       string `json:"sv_enforceGameBuild,omitempty"`
	Maxplayers      string `json:"sv_maxClients,omitempty"`
	Sv_projectName  string `json:"sv_projectName,omitempty"`
	Onesync_enabled string `json:"onesync_enabled,omitempty"`
	Sv_projectDesc  string `json:"sv_projectDesc,omitempty"`
	Tags            string `json:"tags,omitempty"`
}

type fivemplayer struct {
	ID          int      `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Ping        int      `json:"ping,omitempty"`
	Identifiers []string `json:"identifiers,omitempty"`
}

type fivemplayers []fivemplayer

func Fivemquery(hostname string) (string, error) {
	gamedata := FivemQueryAsStruct(hostname)
	djson := jsonify(gamedata)
	return (string(djson)), nil
}

func FivemQueryAsStruct(hostname string) types.Gamequery {
	var gamedata types.Gamequery

	req, err := http.NewRequest("GET", "http://"+hostname+"/info.json", nil)
	if err != nil {
		gamedata.Error = "Server is offline."
		return gamedata
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Warn(err)
		gamedata.Error = "Error"
		return gamedata
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var fiveminfojson1 fiveminfojson
	json.Unmarshal([]byte(string(body)), &fiveminfojson1)

	gamedata.Maxplayers, _ = strconv.Atoi(fiveminfojson1.Vars.Maxplayers)
	gamedata.ServerInfo.Name = fiveminfojson1.Vars.Sv_projectName
	gamedata.ServerInfo.MaxPlayers, _ = strconv.Atoi(fiveminfojson1.Vars.Maxplayers)
	gamedata.ServerInfo.Extra = map[string]string{
		"icon":           fiveminfojson1.Icon,
		"Resources":      strconv.Itoa(len(fiveminfojson1.Resources)),
		"Version":        strconv.Itoa(fiveminfojson1.Version),
		"Server":         fiveminfojson1.Server,
		"Onesync":        fiveminfojson1.Vars.Onesync_enabled,
		"Gamebuild":      fiveminfojson1.Vars.GameBuild,
		"sv_projectDesc": fiveminfojson1.Vars.Sv_projectDesc,
		"tags":           fiveminfojson1.Vars.Tags,
	}

	req, err = http.NewRequest("GET", "http://"+hostname+"/players.json", nil)
	if err != nil {
		gamedata.Error = "Error, unable to query for players"
		return gamedata
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		logrus.Warn(err)
	}
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	var fivemplayers fivemplayers
	json.Unmarshal([]byte(string(body)), &fivemplayers)

	gamedata.Onlineplayers = len(fivemplayers)

	for _, s := range fivemplayers {
		identifier := ""
		if len(s.Identifiers) > 0 {
			identifier = s.Identifiers[0]
		}
		gamedata.Detailed = append(gamedata.Detailed, map[string]string{
			"id":          strconv.Itoa(s.ID),
			"name":        s.Name,
			"identifier":  identifier,
			"identifiers": GetIdentifiersAsString(s),
			"ping":        strconv.Itoa(s.Ping),
		})
	}

	gamedata.ServerInfo.OnlinePlayers = len(fivemplayers)
	gamedata.Error = ""
	return gamedata
}

func GetIdentifiersAsString(s fivemplayer) string {
	resultant := ""
	for _, k := range s.Identifiers {
		resultant += k + "\n"
	}
	return resultant
}
