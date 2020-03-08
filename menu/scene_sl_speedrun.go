package menu

import (
	"os"
	"path/filepath"

	"github.com/libretro/ludo/core"
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/savestates"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/speedrun"
	"github.com/libretro/ludo/state"
)

func buildSLSpeedrunMenu(game speedrun.GameCfg) Scene {
	var list sceneMain
	list.label = "SpeedLearn Speedrun Menu"

	//	rdbGame := cfg.RomInfo

	for _, sr := range game.Speedruns {
		sr := sr
		list.children = append(list.children, entry{
			label:    sr.Name,
			subLabel: sr.Description,
			icon:     "speedrun",
			callbackOK: func() {
				loadSpeedrun(&list, game, sr)
			},
		})
	}

	list.segueMount()

	return &list
}

func loadSpeedrun(scene Scene, game speedrun.GameCfg, sr speedrun.SpeedrunCfg) {
	// Load game core
	if _, err := os.Stat(game.RomPath()); os.IsNotExist(err) {
		ntf.DisplayAndLog(ntf.Error, speedrun.NotificationPrefix, "Game not found: %s", game.RomPath())
		return
	}
	corePath, err := settings.CoreForPlaylist(game.RomInfo.System)
	if err != nil {
		ntf.DisplayAndLog(ntf.Error, speedrun.NotificationPrefix, "Cannot load core: %v", err.Error())
		return
	}
	if _, err := os.Stat(corePath); os.IsNotExist(err) {
		ntf.DisplayAndLog(ntf.Error, speedrun.NotificationPrefix, "Core not found: %s", filepath.Base(corePath))
		return
	}
	if state.Global.CorePath != corePath {
		if err := core.Load(corePath); err != nil {
			ntf.DisplayAndLog(ntf.Error, speedrun.NotificationPrefix, "Cannot load core: %v", err.Error())
			return
		}
	}

	// Load game
	if state.Global.GamePath != game.RomPath() {
		if err := core.LoadGame(game.RomPath()); err != nil {
			ntf.DisplayAndLog(ntf.Error, speedrun.NotificationPrefix, "Cannot load game: %v", err.Error())
			return
		}
	}

	// Load savestate
	err = savestates.Load(sr.SavestatePath())
	if err != nil {
		ntf.DisplayAndLog(ntf.Error, speedrun.NotificationPrefix, err.Error())
		return
	}

	scene.segueNext()
	menu.Push(buildQuickMenu())
	menu.tweens.FastForward() // position the elements without animating
	state.Global.MenuActive = false

	state.Global.SpeedrunSession.Stopwatch.Start()
}
