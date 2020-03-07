package menu

import (
	"github.com/libretro/ludo/speedrun"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"
)

func buildSLGameMenu(cfg []speedrun.GameCfg) (Scene, error) {
	var list scenePlaylist
	list.label = "Speedrun Games Menu"

	for _, c := range cfg {
		game, err := c.RomInfo(state.Global.DB)
		if err != nil {
			return nil, err
		}
		strippedName, tags := extractTags(game.Name)
		list.children = append(list.children, entry{
			label:    strippedName,
			gameName: game.Name,
			path:     game.Path,
			system:   game.System,
			tags:     tags,
			icon:     utils.FileName(c.RomPath()) + "-content",
			callbackOK: func() {
				// TODO:
			},
		})
	}

	list.segueMount()

	return &list, nil
}
