package menu

import (
	"github.com/libretro/ludo/speedrun"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/video"
)

// RenderSpeedrunOverlay draws the current speedrun overlay:
// * Stopwatch at top center of the screen when running
func (m *Menu) RenderSpeedrunOverlay() {
	fbw, fbh := vid.Window.GetFramebufferSize()
	vid.Font.UpdateResolution(fbw, fbh)

	sw := state.Global.SpeedrunSession.Stopwatch
	swString := speedrun.StopwatchFormatMinute(sw.Elapsed())

	if state.Global.CoreRunning && !state.Global.MenuActive {
		// White font and black transparent background as menu
		lw := vid.Font.Width(m.ratio, swString)
		rectW := lw + 30*m.ratio
		rectH := 120 * menu.ratio

		rectColor := video.Color{R: 0, G: 0, B: 0, A: 0.5}
		vid.DrawRect(
			(float32(fbw)-rectW)/2, -rectH*0.2, rectW, rectH, 0.2, rectColor,
		)

		vid.Font.SetColor(1, 1, 1, 1)
		vid.Font.Printf(
			(float32(fbw)-lw)/2,
			75*m.ratio,
			m.ratio,
			swString,
		)
	}
}
