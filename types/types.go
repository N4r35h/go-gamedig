package types

type ServerInfo struct {
	Name          string
	OnlinePlayers int
	MaxPlayers    int
	Extra         map[string]string
}

type Gamequery struct {
	ServerInfo    ServerInfo
	Error         string              `json:"error"`
	Onlineplayers int                 `json:"onlineplayers"`
	Maxplayers    int                 `json:"maxplayers"`
	Detailed      []map[string]string `json:"detailed"`
	CurrentTime   int64               `json:"currenttime"`
}
