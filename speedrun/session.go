// Package speedrun manage speedrun in Ludo
package speedrun

// Session manage all information and action relative to a speedrun session
type Session struct {
	Player      string
	Stopwatch   Stopwatch
	CloudDB     *CloudDB
	GameCfg     GameCfg
	CategoryCfg CategoryCfg
}
