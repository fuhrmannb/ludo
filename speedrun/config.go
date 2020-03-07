package speedrun

import (
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"path/filepath"

	"github.com/libretro/ludo/rdb"
	"github.com/libretro/ludo/utils"
	"gopkg.in/yaml.v2"
)

// NotificationPrefix is used by notification system
var NotificationPrefix = "Speedrun"

// GameCfg is the game configuration containing all speedrun information
type GameCfg struct {
	Rom       string        `yaml:"rom"`
	Speedruns []SpeedrunCfg `yaml:"speedruns"`

	speedrunDir string
}

// SpeedrunCfg contains a speedrun setup configuration for a specific game
type SpeedrunCfg struct {
	Name        string    `yaml:"name"`
	Savestate   string    `yaml:"savestate"`
	Description string    `yaml:"description"`
	StartSplit  string    `yaml:"start_split"`
	EndSplit    string    `yaml:"end_split"`
	Difficulty  int       `yaml:"difficulty"`
	Trophies    TrophyCfg `yaml:"trophies"`
}

// TrophyCfg contains the list of trophies time for a speedrun
// You earn the relative trophy when you beat the specified time
type TrophyCfg struct {
	Copper   string `yaml:"copper"`
	Silver   string `yaml:"silver"`
	Gold     string `yaml:"gold"`
	Platinum string `yaml:"platinum"`
}

// LoadCfg load the speedrun config from the Speedrun directory
func LoadCfg(speedrunDir string) ([]GameCfg, error) {

	tomlFiles, err := filepath.Glob(fmt.Sprintf("%s/*.yaml", speedrunDir))
	if err != nil {
		return nil, fmt.Errorf("cannot find speedrun YAML files: %v", err)
	}

	var games []GameCfg

	for _, tf := range tomlFiles {
		b, err := ioutil.ReadFile(tf)
		if err != nil {
			return nil, fmt.Errorf("cannot open speedrun YAML file %v: %v", tf, err)
		}
		var cfg GameCfg
		cfg.speedrunDir = speedrunDir
		err = yaml.Unmarshal(b, &cfg)
		if err != nil {
			return nil, fmt.Errorf("cannot unmarshal speedrun YAML file %v: %v", tf, err)
		}
		games = append(games, cfg)
	}

	return games, nil
}

// RomInfo returns game information comparing RDB database
func (cfg *GameCfg) RomInfo(db rdb.DB) (rdb.Game, error) {
	games := make(chan (rdb.Game))

	// Load game from cfg
	romPath := cfg.RomPath()
	bytes, err := ioutil.ReadFile(romPath)
	if err != nil {
		return rdb.Game{}, fmt.Errorf("cannot load speedrun ROM file %v: %v", cfg.Rom, err)
	}

	// Find game in database
	var game rdb.Game
	found := false
	go func() {
		for g := range games {
			found = true
			game = g
		}
		close(games)
	}()

	crc := crc32.ChecksumIEEE(bytes)
	db.FindByCRC(romPath, utils.FileName(romPath), crc, games)

	if !found {
		return rdb.Game{}, fmt.Errorf("game not found in RDB")
	}

	return game, nil
}

// RomPath returns the absolute rom file path
func (cfg *GameCfg) RomPath() string {
	return filepath.Join(cfg.speedrunDir, cfg.Rom)
}
