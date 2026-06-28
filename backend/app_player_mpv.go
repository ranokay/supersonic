//go:build mpv

package backend

import (
	"fmt"

	"github.com/dweymouth/supersonic/backend/player/mpv"
)

// initLocalPlayer constructs and initializes the mpv-backed local player.
func (a *App) initLocalPlayer() error {
	p := mpv.NewWithClientName(a.appName)
	c := a.Config.LocalPlayback
	c.InMemoryCacheSizeMB = clamp(c.InMemoryCacheSizeMB, 10, 500)
	if err := p.Init(c.InMemoryCacheSizeMB); err != nil {
		return fmt.Errorf("failed to initialize mpv player: %s", err.Error())
	}
	a.LocalPlayer = p
	return nil
}
