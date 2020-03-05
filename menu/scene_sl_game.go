package menu

type sceneSLGame struct {
	entry
}

func buildSLGameMenu() Scene {
	var list sceneSLGame
	list.label = "Speedrun game"

	//TODO: game menu

	list.segueMount()

	return &list
}

func (s *sceneSLGame) Entry() *entry {
	return &s.entry
}

func (s *sceneSLGame) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *sceneSLGame) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *sceneSLGame) segueBack() {
	genericAnimate(&s.entry)
}

func (s *sceneSLGame) update(dt float32) {
	genericInput(&s.entry, dt)
}

func (s *sceneSLGame) render() {
	genericRender(&s.entry)
}

func (s *sceneSLGame) drawHintBar() {
	genericDrawHintBar()
}
