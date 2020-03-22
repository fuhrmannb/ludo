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
	ID         string        `yaml:id`
	Rom        string        `yaml:"rom"`
	Categories []CategoryCfg `yaml:"categories"`
	RomInfo    rdb.Game      `yaml:-`

	speedrunDir string `yaml:-`
}

// CategoryCfg contains a speedrun setup configuration for a specific game
type CategoryCfg struct {
	ID          string    `yaml:id`
	Name        string    `yaml:"name"`
	Savestate   string    `yaml:"savestate"`
	Description string    `yaml:"description"`
	StartSplit  string    `yaml:"start_split"`
	EndSplit    string    `yaml:"end_split"`
	Difficulty  int       `yaml:"difficulty"`
	Trophies    TrophyCfg `yaml:"trophies"`

	speedrunDir string `yaml:-`
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
func LoadCfg(speedrunDir string, db rdb.DB) ([]GameCfg, error) {

	tomlFiles, err := filepath.Glob(fmt.Sprintf("%s/*.yaml", speedrunDir))
	if err != nil {
		return nil, fmt.Errorf("cannot find speedrun YAML files: %w", err)
	}

	var games []GameCfg

	for _, tf := range tomlFiles {
		b, err := ioutil.ReadFile(tf)
		if err != nil {
			return nil, fmt.Errorf("cannot open speedrun YAML file %v: %w", tf, err)
		}
		var cfg GameCfg
		err = yaml.Unmarshal(b, &cfg)
		if err != nil {
			return nil, fmt.Errorf("cannot unmarshal speedrun YAML file %v: %w", tf, err)
		}

		cfg.speedrunDir = speedrunDir
		for i := range cfg.Categories {
			cfg.Categories[i].speedrunDir = speedrunDir
		}

		cfg.RomInfo, err = romInfo(cfg, db)
		if err != nil {
			return nil, fmt.Errorf("cannot get rom info from %v: %w", tf, err)
		}

		games = append(games, cfg)
	}

	return games, nil
}

func romInfo(cfg GameCfg, db rdb.DB) (rdb.Game, error) {
	games := make(chan (rdb.Game))

	// Load game from cfg
	romPath := cfg.RomPath()
	bytes, err := ioutil.ReadFile(romPath)
	if err != nil {
		return rdb.Game{}, fmt.Errorf("cannot load speedrun ROM file %v: %w", cfg.Rom, err)
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
func (gc *GameCfg) RomPath() string {
	return filepath.Join(gc.speedrunDir, gc.Rom)
}

// SavestatePath returns the absolute stavestate path
func (sc *CategoryCfg) SavestatePath() string {
	return filepath.Join(sc.speedrunDir, sc.Savestate)
}
