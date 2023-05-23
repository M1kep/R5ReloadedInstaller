package page

import (
	"R5ReloadedInstaller/internal/ui/interaction"
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// checkboxes is a list of all checkboxes
var checkboxes = []widget.Bool{
	installSDKToggle,
	installLatestFLowToggle,
	cleanScriptsToggle,
	installLatestR5ScriptsToggle,
	installPreReleaseSDKToggle,
}

var (
	// Checkboxes
	installSDKToggle             widget.Bool
	installLatestFLowToggle      widget.Bool
	cleanScriptsToggle           widget.Bool
	installLatestR5ScriptsToggle widget.Bool
	installPreReleaseSDKToggle   widget.Bool

	// Buttons
	startButton widget.Clickable

	// checkBoxList is a list of all checkboxes
	checkBoxList []layout.FlexChild
)

type MainPage struct {
	interactionManager *interaction.Manager
}

func (m *MainPage) Name() string {
	return "Main Page"
}

func (m *MainPage) Layout(gtx layout.Context) layout.Dimensions {
	rigidStart := layout.Rigid(
		func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(theme, &startButton, "Start")
			return btn.Layout(gtx)
		})

	rigidEnd := layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout)

	combined := append(checkBoxList, append([]layout.FlexChild{rigidStart}, rigidEnd)...)

	return layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceStart,
	}.Layout(gtx,
		combined...,
	)
}

func checkboxHandler(msg string) func() {
	return func() {
		fmt.Println(msg)
	}
}

func (m *MainPage) OnPageLoad() {
	m.interactionManager = interaction.NewManager()
	// Buttons
	m.interactionManager.RegisterClickableHandler(&startButton, func() {
		fmt.Println("startButtonClicked")
	})

	// Checkboxes
	m.interactionManager.RegisterBoolHandler(&installSDKToggle, checkboxHandler("installSDKToggle changed"))
	m.interactionManager.RegisterBoolHandler(&installLatestFLowToggle, checkboxHandler("installLatestFLowToggle changed"))
	m.interactionManager.RegisterBoolHandler(&cleanScriptsToggle, checkboxHandler("cleanScriptsToggle changed"))
	m.interactionManager.RegisterBoolHandler(&installLatestR5ScriptsToggle, checkboxHandler("installLatestR5ScriptsToggle changed"))
	m.interactionManager.RegisterBoolHandler(&installPreReleaseSDKToggle, checkboxHandler("installPreReleaseSDKToggle changed"))

	checkBoxList = checkBoxListBuilder([]map[*widget.Bool]string{
		{&installSDKToggle: "Install SDK"},
		{&installLatestFLowToggle: "Install latest Flow"},
		{&cleanScriptsToggle: "Clean scripts"},
		{&installLatestR5ScriptsToggle: "Install latest R5 scripts"},
		{&installPreReleaseSDKToggle: "Install pre-release SDK"},
	})
}

func (m *MainPage) OnPageUnload() {
	return
}

func (m *MainPage) HandleUserInteractions() {
	m.interactionManager.RunHandlers()
}
