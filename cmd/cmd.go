package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	gamequery "github.com/N4r35h/go-gamedig/games"
	"github.com/N4r35h/go-gamedig/types"
	"github.com/spf13/cobra"
)

var GameQueryHandlers map[string]func(hostname string) types.Gamequery = map[string]func(hostname string) types.Gamequery{
	"samp":           gamequery.SampqueryAsStruct,
	"fivem":          gamequery.FivemQueryAsStruct,
	"minecraft":      gamequery.MinecraftqueryAsStruct,
	"ark":            gamequery.ValveQueryAsStruct,
	"spaceengineers": gamequery.ValveQueryAsStruct,
	"7d2d":           gamequery.SevenD2DQueryAsStruct,
	"tf2":            gamequery.SevenD2DQueryAsStruct,
	"unturned":       gamequery.SevenD2DQueryAsStruct,
	"csgo":           gamequery.CSGOQueryAsStruct,
	"mta":            gamequery.MtaqueryAsStruct,
	"vcmp":           gamequery.VCMPqueryAsStruct,
	"gtac":           gamequery.GTACqueryAsStruct,
}

var rootCmd = &cobra.Command{
	Use:   "gogamedig [game_type] [endpoint]",
	Short: "Go Game Dig is a CLI Utility to query gameservers",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("requires at least two arg")
		}
		if _, exists := GameQueryHandlers[args[0]]; !exists {
			return errors.New("unsupported game_type")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		op := GameQueryHandlers[args[0]](args[1])
		marshaled, err := json.MarshalIndent(op, "", "   ")
		if err != nil {
			log.Fatalf("marshaling error: %s", err)
		}
		fmt.Println(string(marshaled))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
