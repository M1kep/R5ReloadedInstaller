package page

import (
	"R5ReloadedInstaller/internal/app"
	"R5ReloadedInstaller/internal/ui/interaction"
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type TroubleshootingPage struct {
	interactionManager *interaction.Manager
	windowManager      app.WindowManager

	// Buttons
	runButton  widget.Clickable
	checkboxes map[string]*widget.Bool
}

func (t *TroubleshootingPage) ID() string {
	return "Troubleshooting Page"
}

func (t *TroubleshootingPage) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, material.H1(theme, "Troubleshooting").Layout)
			},
		),
		layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, material.Body1(theme, "Actions on this page should only be ran if you are sure about what they are doing.").Layout)
			},
		),
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				chkBox := material.CheckBox(theme, t.checkboxes["Clean Scripts"], "Clean Scripts")
				return chkBox.Layout(gtx)
			},
		),
		layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(theme, &t.runButton, "Run")
				return btn.Layout(gtx)
			},
		),
		layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(theme, &backButton, "Back")
				return btn.Layout(gtx)
			},
		),
		layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
	)
}

func (t *TroubleshootingPage) OnPageLoad(windowManager app.WindowManager) {
	if t.interactionManager != nil {
		return
	}

	t.interactionManager = interaction.NewManager()
	t.windowManager = windowManager

	t.interactionManager.RegisterClickableHandler(&backButton, func() {
		if t.windowManager != nil {
			t.windowManager.PreviousPage()
		}
	})

	t.interactionManager.RegisterClickableHandler(&t.runButton, func() {
		// Print values of checkboxes
		for name, checkbox := range t.checkboxes {
			fmt.Printf("%s: %t\n", name, checkbox.Value)
		}
	})

	cleanScriptsToggle := widget.Bool{}
	t.checkboxes = map[string]*widget.Bool{
		"Clean Scripts": &cleanScriptsToggle,
	}
}

func (t *TroubleshootingPage) OnPageUnload() {
	return
}

func (t *TroubleshootingPage) HandleUserInteractions() {
	t.interactionManager.RunHandlers()
}
