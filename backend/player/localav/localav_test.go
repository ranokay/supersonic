//go:build !mpv && darwin && integration

package localav_test

import (
	"testing"
	"time"

	"github.com/dweymouth/supersonic/backend/mediaprovider"
	"github.com/dweymouth/supersonic/backend/player"
	"github.com/dweymouth/supersonic/backend/player/localav"
)

func TestPlayerBasic(t *testing.T) {
	p := localav.New()
	if err := p.Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	defer p.Destroy()

	stopped := make(chan struct{}, 1)
	p.OnStopped(func() { stopped <- struct{}{} })

	err := p.PlayFile("/System/Library/Sounds/Ping.aiff", mediaprovider.MediaItemMetadata{}, 0)
	if err != nil {
		t.Fatalf("PlayFile: %v", err)
	}

	waitForState(t, p, player.Playing, time.Second)

	// Wait for natural end (Ping is ~0.5s)
	select {
	case <-stopped:
		t.Log("stopped naturally")
	case <-time.After(5 * time.Second):
		t.Error("timed out waiting for stop")
	}
}

func TestPlayerPauseResume(t *testing.T) {
	p := localav.New()
	if err := p.Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	defer p.Destroy()

	err := p.PlayFile("/System/Library/Sounds/Submarine.aiff", mediaprovider.MediaItemMetadata{}, 0)
	if err != nil {
		t.Fatalf("PlayFile: %v", err)
	}
	waitForState(t, p, player.Playing, time.Second)

	if err := p.Pause(); err != nil {
		t.Fatalf("Pause: %v", err)
	}
	waitForState(t, p, player.Paused, time.Second)

	if err := p.Continue(); err != nil {
		t.Fatalf("Continue: %v", err)
	}
	waitForState(t, p, player.Playing, time.Second)

	p.Stop(false)
}

func TestPlayerSeek(t *testing.T) {
	p := localav.New()
	if err := p.Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	defer p.Destroy()

	// Submarine is ~1.8s — seek to 1s
	err := p.PlayFile("/System/Library/Sounds/Submarine.aiff", mediaprovider.MediaItemMetadata{}, 0)
	if err != nil {
		t.Fatalf("PlayFile: %v", err)
	}
	waitForState(t, p, player.Playing, time.Second)

	if err := p.SeekSeconds(1.0); err != nil {
		t.Fatalf("SeekSeconds: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
	pos := p.GetStatus().TimePos
	if pos < 0.5 || pos > 1.8 {
		t.Errorf("position after seek = %v, expected ~1.0", pos)
	}
	p.Stop(false)
}

func waitForState(t *testing.T, p *localav.Player, want player.State, timeout time.Duration) {
	t.Helper()
	deadline := time.After(timeout)
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		if got := p.GetStatus().State; got == want {
			return
		}
		select {
		case <-deadline:
			t.Fatalf("timed out waiting for state %v; got %v", want, p.GetStatus().State)
		case <-ticker.C:
		}
	}
}

func TestPlayerListDevices(t *testing.T) {
	p := localav.New()
	if err := p.Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	defer p.Destroy()

	devices, err := p.ListAudioDevices()
	if err != nil {
		t.Fatalf("ListAudioDevices: %v", err)
	}
	if len(devices) == 0 {
		t.Error("expected at least one audio device")
	}
	t.Logf("devices: %+v", devices)
}

func TestPlayerVolume(t *testing.T) {
	p := localav.New()
	if err := p.Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	defer p.Destroy()

	p.SetVolume(50)
	if v := p.GetVolume(); v != 50 {
		t.Errorf("expected volume 50, got %v", v)
	}
}
