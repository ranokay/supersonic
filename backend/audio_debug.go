//go:build debug

package backend

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type audioDebugSpan struct {
	label string
	start time.Time
}

func startAudioDebugSpan(label string, fields ...any) audioDebugSpan {
	log.Printf("[audio-start] %s begin%s", label, audioDebugFields(fields...))
	return audioDebugSpan{label: label, start: time.Now()}
}

func (s audioDebugSpan) Done(fields ...any) {
	log.Printf("[audio-start] %s done duration=%s%s", s.label, time.Since(s.start).Round(time.Millisecond), audioDebugFields(fields...))
}

func audioDebugf(format string, args ...any) {
	log.Printf("[audio-start] "+format, args...)
}

func audioDebugFields(fields ...any) string {
	if len(fields) == 0 {
		return ""
	}
	var b strings.Builder
	for i := 0; i+1 < len(fields); i += 2 {
		b.WriteByte(' ')
		b.WriteString(fields[i].(string))
		b.WriteByte('=')
		b.WriteString(toDebugString(fields[i+1]))
	}
	return b.String()
}

func toDebugString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case bool:
		if t {
			return "true"
		}
		return "false"
	default:
		return strings.TrimSpace(fmt.Sprint(v))
	}
}
