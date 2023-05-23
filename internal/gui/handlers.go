package gui

import (
	"fmt"
	"gioui.org/widget"
)

var buttonHandlers = make(map[*widget.Clickable]func())
var boolHandlers = make(map[*widget.Bool]func())

func runHandlers() {
	// Checkboxes come first as they may change values referenced by other handlers
	for bool, handler := range boolHandlers {
		if bool.Changed() {
			handler()
		}
	}

	for clickable, handler := range buttonHandlers {
		if clickable.Clicked() {
			handler()
		}
	}
}

func registerClickableHandler(clickable *widget.Clickable, handler func()) {
	buttonHandlers[clickable] = handler
}

func registerBoolHandler(bool *widget.Bool, handler func()) {
	boolHandlers[bool] = handler
}

func startButtonClicked() {
	fmt.Println("startButtonClicked")
}
