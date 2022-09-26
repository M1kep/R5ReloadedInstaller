package main

import (
	"R5ReloadedInstaller/pkg/util"
	"R5ReloadedInstaller/pkg/validation"
	"fmt"
	"github.com/google/go-github/v47/github"
	"github.com/gosuri/uiprogress"
	"github.com/pkg/browser"
	"github.com/rs/zerolog"
	"github.com/tawesoft/golib/v2/dialog"
	"golang.org/x/sync/errgroup"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	VERSION := "v0.13.0"
	var r5Folder string
	ghClient := github.NewClient(nil)

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
	//options = append(options, optionConfig{"Aim Trainer", 100, 300})
	//options = append(options, optionConfig{
	//	"(Troubleshooting) Clean Scripts - Deletes 'platform/scripts' prior to extracting",
	//	1000,
	//	50,
	//})
	selectedOptions, err := gatherRunOptions([]string{
		"SDK",
		"Aim Trainer",
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
	errGroup := new(errgroup.Group)

	if util.Contains(selectedOptions, "(Troubleshooting) Clean Scripts - Deletes 'platform/scripts' prior to extracting") {
		err := os.RemoveAll(filepath.Join(r5Folder, "platform/scripts"))
		if err != nil {
			fileLogger.Error().Err(fmt.Errorf("error removing 'platform/scripts' folder: %v", err)).Msg("error")
			util.LogErrorWithDialog(fmt.Errorf("error removing 'platform/scripts' folder: %v", err))
			return
		}
	}

	if util.Contains(selectedOptions, "SDK") {
		err := ProcessSDK(
			ghClient,
			errGroup,
			cacheDir,
			r5Folder,
			false,
		)

		if err != nil {
			fileLogger.Error().Err(err).Msg("error")
			util.LogErrorWithDialog(err)
			return
		}
	}

	if util.Contains(selectedOptions, "SDK(Include Pre-Releases)") {
		err := ProcessSDK(
			ghClient,
			errGroup,
			cacheDir,
			r5Folder,
			true,
		)

		if err != nil {
			fileLogger.Error().Err(err).Msg("error")
			util.LogErrorWithDialog(err)
			return
		}
	}

	if util.Contains(selectedOptions, "(DEV) Latest r5_scripts") {
		err := ProcessLatestR5Scripts(
			ghClient,
			errGroup,
			cacheDir,
			r5Folder,
		)

		if err != nil {
			fileLogger.Error().Err(err).Msg("error")
			util.LogErrorWithDialog(err)
			return
		}
	}

	if util.Contains(selectedOptions, "Aim Trainer") {
		err := ProcessAimTrainer(
			ghClient,
			errGroup,
			cacheDir,
			r5Folder,
		)

		if err != nil {
			fileLogger.Error().Err(err).Msg("error")
			util.LogErrorWithDialog(err)
			return
		}
	}

	_ = dialog.Raise("Success. Confirm to close terminal.")
	return
}
