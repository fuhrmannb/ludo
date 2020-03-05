package menu

import (
	"github.com/libretro/ludo/state"
)

type sceneSLMain struct {
	entry
}

func buildSLMainMenu() Scene {
	var list sceneSLMain
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
			list.segueNext()
			menu.Push(buildSLGameMenu())
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

func (s *sceneSLMain) Entry() *entry {
	return &s.entry
}

func (s *sceneSLMain) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *sceneSLMain) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *sceneSLMain) segueBack() {
	genericAnimate(&s.entry)
}

func (s *sceneSLMain) update(dt float32) {
	genericInput(&s.entry, dt)
}

func (s *sceneSLMain) render() {
	genericRender(&s.entry)
}

func (s *sceneSLMain) drawHintBar() {
	genericDrawHintBar()
}
