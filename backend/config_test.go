package backend

import "testing"

func TestLocalPlaybackStabilizationDefaults(t *testing.T) {
	cfg := DefaultConfig("test")

	if cfg.LocalPlayback.DACWarmUpEnabled {
		t.Fatal("DAC warm-up should default off")
	}
	if got := cfg.LocalPlayback.DACWarmUpDurationSeconds; got != 3 {
		t.Fatalf("DAC warm-up duration = %d, want 3", got)
	}
	if !cfg.LocalPlayback.SampleRateSwitchPauseEnabled {
		t.Fatal("sample-rate switch pause should default on")
	}
	if got := cfg.LocalPlayback.SampleRateSwitchPauseMilliseconds; got != 200 {
		t.Fatalf("sample-rate switch pause = %d, want 200", got)
	}
}

func TestLocalPlaybackStabilizationOptionsClamp(t *testing.T) {
	opts := LocalPlaybackStabilizationOptions(LocalPlaybackConfig{
		DACWarmUpEnabled:                  true,
		DACWarmUpDurationSeconds:          99,
		SampleRateSwitchPauseEnabled:      true,
		SampleRateSwitchPauseMilliseconds: -10,
	})

	if !opts.DACWarmUpEnabled || opts.DACWarmUpDurationSeconds != 10 {
		t.Fatalf("warm-up options = %+v, want enabled with 10s clamp", opts)
	}
	if !opts.SampleRateSwitchPauseEnabled || opts.SampleRateSwitchPauseMilliseconds != 0 {
		t.Fatalf("sample-rate pause options = %+v, want enabled with 0ms clamp", opts)
	}
}
