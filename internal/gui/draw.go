package gui

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var ops op.Ops

// startButton is a clickable widget
var startButton widget.Clickable

// Checkboxes
var (
	installSDKToggle             widget.Bool
	installLatestFLowToggle      widget.Bool
	cleanScriptsToggle           widget.Bool
	installLatestR5ScriptsToggle widget.Bool
	installPreReleaseSDKToggle   widget.Bool
)

// checkboxes is a list of all checkboxes
var checkboxes = []widget.Bool{
	installSDKToggle,
	installLatestFLowToggle,
	cleanScriptsToggle,
	installLatestR5ScriptsToggle,
	installPreReleaseSDKToggle,
}

var checkBoxList []layout.FlexChild

var th = material.NewTheme(gofont.Collection())

func handleFrameEvent(e system.FrameEvent) {
	// Runs any input handlers
	runHandlers()
	// FrameEvent is sent when the window should be redrawn.
	gtx := layout.NewContext(&ops, e)

	rigidStart := layout.Rigid(
		func(gtx layout.Context) D {
			btn := material.Button(th, &startButton, "Start")
			return btn.Layout(gtx)
		})

	rigidEnd := layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout)

	combined := append(checkBoxList, append([]layout.FlexChild{rigidStart}, rigidEnd)...)

	layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceStart,
	}.Layout(gtx,
		combined...,
	)

	e.Frame(gtx.Ops)
}

func checkboxHandler(msg string) func() {
	return func() {
		fmt.Println(msg)
	}
}

func Draw(w *app.Window) error {
	//page := page2.MainPage{}

	//registerClickableHandler(&startButton, startButtonClicked)

	//checkBoxList = checkBoxListBuilder([]map[*widget.Bool]string{
	//	{&installSDKToggle: "Install SDK"},
	//	{&installLatestFLowToggle: "Install latest Flow"},
	//	{&cleanScriptsToggle: "Clean scripts"},
	//	{&installLatestR5ScriptsToggle: "Install latest R5 scripts"},
	//	{&installPreReleaseSDKToggle: "Install pre-release SDK"},
	//})

	//for i, checkBox := range checkboxes {
	//	registerBoolHandler(&checkBox, checkboxHandler(fmt.Sprintf("Checkbox %d changed", i)))
	//}

	registerBoolHandler(&installSDKToggle, checkboxHandler("installSDKToggle changed"))
	registerBoolHandler(&installLatestFLowToggle, checkboxHandler("installLatestFLowToggle changed"))
	registerBoolHandler(&cleanScriptsToggle, checkboxHandler("cleanScriptsToggle changed"))
	registerBoolHandler(&installLatestR5ScriptsToggle, checkboxHandler("installLatestR5ScriptsToggle changed"))
	registerBoolHandler(&installPreReleaseSDKToggle, checkboxHandler("installPreReleaseSDKToggle changed"))

	//runtime.Breakpoint()
	for e := range w.Events() {
		switch e := e.(type) {

		case system.FrameEvent:
			handleFrameEvent(e)
		case system.DestroyEvent:
			return e.Err
		}
	}

	return nil
}
