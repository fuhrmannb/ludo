package menu

import (
	"github.com/libretro/ludo/speedrun"
	"github.com/libretro/ludo/utils"
)

func buildSLGameMenu(cfg []speedrun.GameCfg) (Scene, error) {
	var list scenePlaylist
	list.label = "Speedrun Games Menu"

	for _, c := range cfg {
		c := c
		game := c.RomInfo
		strippedName, tags := extractTags(game.Name)
		list.children = append(list.children, entry{
			label:    strippedName,
			gameName: game.Name,
			path:     game.Path,
			system:   game.System,
			tags:     tags,
			icon:     utils.FileName(c.RomPath()) + "-content",
			callbackOK: func() {
				list.segueNext()
				menu.Push(buildSLSpeedrunMenu(c))
			},
		})
	}

	list.segueMount()

	return &list, nil
}
