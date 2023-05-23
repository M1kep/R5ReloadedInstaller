package main

import (
	"R5ReloadedInstaller/internal/gui"
	page2 "R5ReloadedInstaller/internal/ui/page"
	"R5ReloadedInstaller/pkg/util"
	"R5ReloadedInstaller/pkg/validation"
	"fmt"
	"gioui.org/app"
	"github.com/google/go-github/v52/github"
	"github.com/gosuri/uiprogress"
	"github.com/pkg/browser"
	"github.com/rs/zerolog"
	"github.com/tawesoft/golib/v2/dialog"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//func desktopLayout(version string) *fyne.Container {
//	title := canvas.NewText("R5Reloaded Installer", color.RGBA{A: 255})
//	title.Alignment = fyne.TextAlignCenter
//	title.TextStyle.Bold = true
//	title.TextSize = 30
//
//	return container.NewGridWithRows(3,
//		title,
//		layout.NewSpacer(),
//		canvas.NewText("Version: "+version, color.RGBA{A: 255}),
//		canvas.NewText("By M1kep", color.RGBA{A: 255}),
//	)
//}

//func guiStartup() {
//	mainPage := page2.MainPage{}
//	windowManager := gui.NewWindowManager("R5Reloaded Installer", &mainPage)
//	if err := windowManager.Run(); err != nil {
//		log.Fatal(err)
//	}

//w := app.NewWindow(
//	app.Title("R5Reloaded Installer"),
//	app.Size(unit.Dp(400), unit.Dp(600)),
//)

//var ops op.Ops
//// startButton is a clickable widget
//var startButton widget.Clickable
//
//for e := range w.Events() {
//	switch e := e.(type) {
//
//	case system.FrameEvent:
//		fmt.Println("FrameEvent", e)
//		// ops are the operations that are sent to the window
//
//		th := material.NewTheme(gofont.Collection())
//		// FrameEvent is sent when the window should be redrawn.
//		gtx := layout.NewContext(&ops, e)
//
//		layout.Flex{
//			Axis:    layout.Vertical,
//			Spacing: layout.SpaceStart,
//		}.Layout(gtx,
//			layout.Rigid(
//				func(gtx layout.Context) layout.Dimensions {
//					btn := material.Button(th, &startButton, "Start")
//					return btn.Layout(gtx)
//				}),
//			layout.Rigid(
//				layout.Spacer{Height: unit.Dp(10)}.Layout,
//			),
//		)
//
//		e.Frame(gtx.Ops)
//		//case system.DestroyEvent:
//		//	return e.Err
//	}
//}

//if err := gui.Draw(w); err != nil {
//	log.Fatal(err)
//}

//os.Exit(0)
//}

func main() {
	VERSION := "v0.15.1"
	var r5Folder string
	ghClient := github.NewClient(nil)

	go func() {
		mainPage := page2.MainPage{}
		windowManager := gui.NewWindowManager("R5Reloaded Installer", &mainPage)
		if err := windowManager.Run(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()

	//myApp := app.New()
	//myWindow := myApp.NewWindow("R5Reloaded Installer")

	//titleText := canvas.NewText("R5Reloaded Installer", color.RGBA{R: 255, A: 255})
	//topTitle := container.New(layout.NewHBoxLayout(), titleText)
	//tabs := container.NewAppTabs(
	//	container.NewTabItem("Installation", canvas.NewRectangle(color.RGBA{R: 255})),
	//	container.NewTabItem("Settings",
	//		container.NewGridWithColumns(2,
	//			container.NewGridWithRows(2,
	//				randomRectangles(2)...,
	//			),
	//			container.NewGridWithRows(2,
	//				randomRectangles(2)...,
	//			),
	//		),
	//	),
	//)

	//myWindow.SetContent(desktopLayout(VERSION))
	//myCanvas := myWindow.Canvas()
	//content := container.New(layout.NewGridLayout(3))
	//content.Add(canvas.NewText("R5Reloaded Installer", color.RGBA{A: 255}))
	//content.Add(canvas.NewText("Version: "+VERSION, color.RGBA{A: 255}))
	//content.Add(canvas.NewText("By M1kep", color.RGBA{A: 255}))
	//myCanvas.SetContent(content)

	//myWindow.Resize(fyne.NewSize(900, 900))
	//myWindow.ShowAndRun()

	r5Folder, err := getValidatedR5Folder()
	if err != nil {
		util.LogErrorWithDialog(err)
		return
	}

	if validation.IsLauncherFileLocked(r5Folder) {
		_ = dialog.Raise("Please close the R5 Launcher before running.")
		return
	}

	cacheDir, err := initializeDirectories(r5Folder)
	if err != nil {
		util.LogErrorWithDialog(err)
		return
	}

	logFile, err := os.Create(filepath.Join(cacheDir, "logfile.txt"))
	if err != nil {
		util.LogErrorWithDialog(fmt.Errorf("error creating logging file: %v", err))
		return
	}
	defer logFile.Close()

	fileLogger := zerolog.New(logFile).With().Logger()

	shouldExit, msg, err := checkForUpdate(ghClient, cacheDir, VERSION)
	if msg != "" {
		fmt.Println(msg)
	}

	if err != nil {
		fmt.Println(err)
	}

	if shouldExit {
		if strings.HasPrefix(msg, "New major version") {
			err := browser.OpenURL("https://github.com/M1kep/R5ReloadedInstaller/releases/latest")
			if err != nil {
				fileLogger.Error().Err(fmt.Errorf("error opening browser to latest release: %v", err)).Msg("error")
				_ = dialog.Error("Error opening browser to latest release. Please manually update from https://github.com/M1kep/R5ReloadedInstaller/releases/latest")
				return
			}
		}
		_ = dialog.Raise("Exiting due to update check.")
		return
	}

	//type optionConfig struct {
	//	UIOption    string
	//	UIPriority  int
	//	RunPriority int
	//}
	//var options []optionConfig
	//options = append(options, optionConfig{"SDK", 50, 300})
	//options = append(options, optionConfig{"Latest Flowstate Scripts", 100, 300})
	//options = append(options, optionConfig{
	//	"(Troubleshooting) Clean Scripts - Deletes 'platform/scripts' prior to extracting",
	//	1000,
	//	50,
	//})
	selectedOptions, err := gatherRunOptions([]string{
		"SDK",
		"Latest Flowstate Scripts",
		"(Troubleshooting) Clean Scripts - Deletes 'platform/scripts' prior to extracting",
		"(DEV) Latest r5_scripts",
		"SDK(Include Pre-Releases)",
	})
	if err != nil {
		fileLogger.Error().Err(fmt.Errorf("error gathering run options")).Msg("error")
		util.LogErrorWithDialog(fmt.Errorf("error gathering run options"))
		return
	}

	uiprogress.Start()
	processManager := ProcessManager{
		ghClient: ghClient,
		errGroup: new(errgroup.Group),
		cacheDir: cacheDir,
		r5Folder: r5Folder,
	}

	if util.Contains(selectedOptions, "(Troubleshooting) Clean Scripts - Deletes 'platform/scripts' prior to extracting") {
		err := os.RemoveAll(filepath.Join(r5Folder, "platform/scripts"))
		if err != nil {
			fileLogger.Error().Err(fmt.Errorf("error removing 'platform/scripts' folder: %v", err)).Msg("error")
			util.LogErrorWithDialog(fmt.Errorf("error removing 'platform/scripts' folder: %v", err))
			return
		}
	}

	if util.Contains(selectedOptions, "SDK") {
		err := processManager.ProcessSDK(false)

		if err != nil {
			fileLogger.Error().Err(err).Msg("error")
			util.LogErrorWithDialog(err)
			return
		}
	}

	if util.Contains(selectedOptions, "SDK(Include Pre-Releases)") {
		err := processManager.ProcessSDK(true)

		if err != nil {
			fileLogger.Error().Err(err).Msg("error")
			util.LogErrorWithDialog(err)
			return
		}
	}

	if util.Contains(selectedOptions, "(DEV) Latest r5_scripts") {
		err := processManager.ProcessLatestR5Scripts()

		if err != nil {
			fileLogger.Error().Err(err).Msg("error")
			util.LogErrorWithDialog(err)
			return
		}
	}

	if util.Contains(selectedOptions, "Latest Flowstate Scripts") {
		err := processManager.ProcessFlowstate()

		if err != nil {
			fileLogger.Error().Err(err).Msg("error")
			util.LogErrorWithDialog(err)
			return
		}
	}

	_ = dialog.Raise("Success. Confirm to close terminal.")
	return
}
