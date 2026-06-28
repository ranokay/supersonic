package widgets

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	fynetooltip "github.com/dweymouth/fyne-tooltip"
)

func TestVolumeControlSoftwareVolumeLock(t *testing.T) {
	test.NewTempApp(t)

	c := NewVolumeControl(35)
	var changed []int
	c.OnSetVolume = func(vol int) {
		changed = append(changed, vol)
	}

	c.SetSoftwareVolumeLocked(true, "locked")
	if !c.SoftwareVolumeLocked() {
		t.Fatal("expected software volume to be locked")
	}
	if !c.icon.Disabled() || !c.slider.Disabled() {
		t.Fatal("expected volume controls to be disabled while locked")
	}
	if got := int(c.slider.Value); got != 100 {
		t.Fatalf("locked volume display = %d, want 100", got)
	}

	c.SetVolume(42)
	if got := int(c.slider.Value); got != 100 {
		t.Fatalf("locked volume display after SetVolume = %d, want 100", got)
	}
	c.onChanged(10)
	if len(changed) != 0 {
		t.Fatalf("locked volume change invoked callback: %v", changed)
	}

	c.SetSoftwareVolumeLocked(false, "")
	if c.SoftwareVolumeLocked() {
		t.Fatal("expected software volume to unlock")
	}
	if c.icon.Disabled() || c.slider.Disabled() {
		t.Fatal("expected volume controls to be enabled after unlock")
	}
	if got := int(c.slider.Value); got != 42 {
		t.Fatalf("unlocked volume display = %d, want latent software volume 42", got)
	}
}

func TestQualityPathPopupStaysOpenAcrossRefreshes(t *testing.T) {
	test.NewTempApp(t)

	a := NewAuxControls(100, false)
	w := test.NewWindow(nil)
	w.SetContent(fynetooltip.AddWindowToolTipLayer(a, w.Canvas()))
	defer w.Close()
	defer fynetooltip.DestroyWindowToolTipLayer(w.Canvas())

	a.SetQualityPath(QualityPathInfo{
		Badge:            "Bit-perfect 96.0 kHz / s32",
		Status:           "Bit-Perfect",
		SourceFormat:     "96.0 kHz / 24-bit / 2ch",
		OutputPath:       "Integer PCM / 96.0 kHz / 2ch",
		BitPerfectActive: true,
	})
	a.showQualityPath()
	if a.qualityPop == nil || !a.qualityPop.Visible() {
		t.Fatal("expected quality popup to be visible")
	}
	original := a.qualityPop

	a.SetQualityPath(QualityPathInfo{
		Badge:            "Bit-perfect 192.0 kHz / s32",
		Status:           "Bit-Perfect",
		SourceFormat:     "192.0 kHz / 24-bit / 2ch",
		OutputPath:       "Integer PCM / 192.0 kHz / 2ch",
		BitPerfectActive: true,
	})

	if a.qualityPop != original {
		t.Fatal("expected quality popup to be updated in place")
	}
	if !a.qualityPop.Visible() {
		t.Fatal("expected quality popup to remain visible after quality refresh")
	}
	popupText := qualityPopupText(a.qualityPop.Content)
	if !strings.Contains(popupText, "192.0 kHz / 24-bit / 2ch") {
		t.Fatalf("expected popup to show refreshed source format, got %q", popupText)
	}
	if !strings.Contains(popupText, "Integer PCM / 192.0 kHz / 2ch") {
		t.Fatalf("expected popup to show refreshed output path, got %q", popupText)
	}
}

func qualityPopupText(obj fyne.CanvasObject) string {
	var parts []string
	collectQualityPopupText(obj, &parts)
	return strings.Join(parts, "\n")
}

func collectQualityPopupText(obj fyne.CanvasObject, parts *[]string) {
	switch o := obj.(type) {
	case *widget.Label:
		*parts = append(*parts, o.Text)
	case *fyne.Container:
		for _, child := range o.Objects {
			collectQualityPopupText(child, parts)
		}
	}
}
