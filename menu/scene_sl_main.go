package menu

import (
	"github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/speedrun"
	"github.com/libretro/ludo/state"
)

func buildSLMainMenu() Scene {
	var list sceneMain
	list.label = "SpeedLearn Main Menu"

	// Speedrun start is available when:
	// * A player has already played (Player state is not empty)
	// * A player has just enter his nickname
	prependEntry := entry{
		label:       "Speedrun as player",
		subLabel:    "Start speedrun session with current player",
		stringValue: func() string { return state.Global.SpeedrunSession.Player },
		icon:        "run",
		callbackOK: func() {
			cfg, err := speedrun.LoadCfg(settings.Current.SpeedrunDirectory, state.Global.DB)
			if err != nil {
				notifications.DisplayAndLog(notifications.Error, speedrun.NotificationPrefix, "%v", err)
				return
			}

			sc, err := buildSLGameMenu(cfg)
			if err != nil {
				notifications.DisplayAndLog(notifications.Error, speedrun.NotificationPrefix, "%v", err)
				return
			}

			list.segueNext()
			menu.Push(sc)
		},
	}
	prependStart := func() {
		if len(list.children) > 0 && list.children[0].label == prependEntry.label {
			return
		}
		list.children = append([]entry{prependEntry}, list.children...)
	}

	if state.Global.SpeedrunSession.Player != "" {
		prependStart()
	}

	list.children = append(list.children, entry{
		label:    "New player...",
		subLabel: "Start a speedrun session",
		icon:     "add",
		callbackOK: func() {
			list.segueNext()
			menu.Push(buildKeyboard(
				"Type your nickname",
				func(nickname string) {
					state.Global.SpeedrunSession.Player = nickname
					prependStart()
				},
			))
		},
	})

	list.children = append(list.children, entry{
		label:    "Continue with existing player...",
		subLabel: "Continue speedrun session with a player that already played",
		icon:     "resume",
		callbackOK: func() {
			// TODO:
		},
	})

	list.segueMount()

	return &list
}

// TODO: improve SL speedrun render (show description, start/end split, trophies...)
