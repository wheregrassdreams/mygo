package ui

import (
	"testing"

	"github.com/egoist/mygo/transfer"
)

func TestHandoffDataDragReleasesSourceAndAllowsNextClick(t *testing.T) {
	tests := []struct {
		name   string
		result transfer.Result
		local  any
	}{
		{name: "drop uses serialized data", result: transfer.Result{Operation: transfer.Move}},
		{name: "cancel retains optional local value", result: transfer.Result{Canceled: true}, local: "local payload"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			startHandoff, handedOff := false, false
			clicks, doneCalls := 0, 0
			data := transfer.TextData("native payload")
			tt := coreNewTester(func(c *context) {
				coreColumn(c).Gap(20).Padding(20).Children(func() {
					source := coreBox(c).Size(100, 40).Label("Source").Drag("local payload")
					coreDrop[string](source)
					target := coreBox(c).Size(120, 40).Label("Next click")
					if target.Clicked() {
						clicks++
					}
					if startHandoff && !handedOff && source.Dragging() {
						handedOff = true
						if err := c.rt.public.HandoffDataDrag(data, test.local, transfer.DragOptions{Done: func(transfer.Result) { doneCalls++ }}); err != nil {
							t.Errorf("handoff: %v", err)
						}
					}
				})
			}, 300, 200)
			tt.Frame()
			source, _ := tt.Find("Source")
			tt.Press(source.X+10, source.Y+10)
			tt.Move(source.X+20, source.Y+10)
			if tt.rt.drag == nil || tt.rt.drag.native {
				t.Fatal("threshold did not start a local drag")
			}
			startHandoff = true
			tt.Move(-10, -10)
			if !handedOff || tt.rt.drag == nil || !tt.rt.drag.native || tt.h.dragLocal != test.local {
				t.Fatalf("handoff did not promote active drag: handed=%v drag=%+v local=%v", handedOff, tt.rt.drag, tt.h.dragLocal)
			}
			if tt.rt.pressed == nil {
				t.Fatal("source press ended before native completion")
			}
			tt.h.finishDataDrag(test.result)
			if tt.rt.drag != nil || tt.rt.pressed != nil || doneCalls != 1 {
				t.Fatalf("native completion did not clean runtime once: drag=%v pressed=%v done=%d", tt.rt.drag, tt.rt.pressed, doneCalls)
			}
			tt.Frame()
			target, ok := tt.Find("Next click")
			if !ok {
				t.Fatal("next click target missing")
			}
			tt.ClickAt(target.X+10, target.Y+10)
			if clicks != 1 {
				t.Fatalf("next click did not work after native completion: %d", clicks)
			}
		})
	}
}
