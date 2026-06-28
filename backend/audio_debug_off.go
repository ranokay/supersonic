//go:build !debug

package backend

type audioDebugSpan struct{}

func startAudioDebugSpan(string, ...any) audioDebugSpan { return audioDebugSpan{} }

func (audioDebugSpan) Done(...any) {}

func audioDebugf(string, ...any) {}
