package gamequery

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

func jsonify(gamedata interface{}) []byte {
	djson, err := json.Marshal(gamedata)
	if err != nil {
		logrus.Warn(err)
	}
	return djson
}
