package main

import (
	"R5ReloadedInstaller/internal/download"
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
	VERSION := "v0.8.0"
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

	selectedOptions, err := gatherRunOptions([]string{
		"SDK",
		"Aim Trainer",
		"(Troubleshooting) Clean Scripts - Deletes 'platform/scripts' prior to extracting",
		"(DEV) Latest r5_scripts",
	})
	if err != nil {
		fileLogger.Error().Err(fmt.Errorf("error gathering run options")).Msg("error")
		util.LogErrorWithDialog(fmt.Errorf("error gathering run options"))
		return
	}

	uiprogress.Start()
	errGroup := new(errgroup.Group)
	// PHASE 1
	if util.Contains(selectedOptions, "(Troubleshooting) Clean Scripts - Deletes 'platform/scripts' prior to extracting") {
		err := os.RemoveAll(filepath.Join(r5Folder, "platform/scripts"))
		if err != nil {
			fileLogger.Error().Err(fmt.Errorf("error removing 'platform/scripts' folder: %v", err)).Msg("error")
			util.LogErrorWithDialog(fmt.Errorf("error removing 'platform/scripts' folder: %v", err))
			return
		}
	}

	if util.Contains(selectedOptions, "SDK") {
		// Download SDK Release
		sdkOutputPath, err := download.StartLatestRepoReleaseDownload(
			ghClient,
			errGroup,
			"Downloading SDK",
			cacheDir,
			"sdk-depot",
			"depot.zip",
			"Mauler125",
			"r5sdk",
		)
		if err != nil {
			fileLogger.Error().Err(fmt.Errorf("error starting download of sdk release: %v", err)).Msg("error")
			util.LogErrorWithDialog(fmt.Errorf("error starting download of sdk release: %v", err))
			return
		}

		if err := errGroup.Wait(); err != nil {
			fileLogger.Error().Err(fmt.Errorf("error encountered while performing SDK download: %v", err)).Msg("error")
			util.LogErrorWithDialog(fmt.Errorf("error encountered while performing SDK download: %v", err))
			return
		}

		// Unzip SDK into R5Folder
		err = util.UnzipFile(sdkOutputPath, r5Folder, false, "Extracting SDK")
		if err != nil {
			fileLogger.Error().Err(fmt.Errorf("error unzipping sdk: %v", err)).Msg("error")
			util.LogErrorWithDialog(fmt.Errorf("error unzipping sdk: %v", err))
			return
		}
	}

	if util.Contains(selectedOptions, "(DEV) Latest r5_scripts") {
		// Download scripts_r5
		scriptsRepoContentsOutput, err := download.StartLatestRepoContentsDownload(
			ghClient,
			errGroup,
			"Downloading Scripts",
			cacheDir,
			"scripts",
			"Mauler125",
			"scripts_r5",
		)
		if err != nil {
			fileLogger.Error().Err(fmt.Errorf("error starting download of scripts: %v", err)).Msg("error")
			util.LogErrorWithDialog(fmt.Errorf("error starting download of scripts: %v", err))
			return
		}

		if err := errGroup.Wait(); err != nil {
			fileLogger.Error().Err(fmt.Errorf("error encountered while performing r5_scripts download: %v", err)).Msg("error")
			util.LogErrorWithDialog(fmt.Errorf("error encountered while performing r5_scripts download: %v", err))
			return
		}

		// Unzip Scripts into platform/scripts
		err = util.UnzipFile(scriptsRepoContentsOutput, filepath.Join(r5Folder, "platform/scripts"), true, "Extracting scripts")
		if err != nil {
			fileLogger.Error().Err(fmt.Errorf("error unzipping scripts: %v", err)).Msg("error")
			util.LogErrorWithDialog(fmt.Errorf("error unzipping scripts: %v", err))
			return
		}
	}

	if util.Contains(selectedOptions, "Aim Trainer") {
		// Download Aim Trainer release
		aimTrainerReleaseOutput, err := download.StartLatestRepoReleaseDownload(
			ghClient,
			errGroup,
			"Downloading Aim Trainer",
			cacheDir,
			"aimtrainer-deps",
			"AimTrainerRequiredFiles.zip",
			"ColombianGuy",
			"r5_aimtrainer",
		)
		if err != nil {
			fileLogger.Error().Err(fmt.Errorf("error starting download of aimtrainer release: %v", err)).Msg("error")
			util.LogErrorWithDialog(fmt.Errorf("error starting download of aimtrainer release: %v", err))
			return
		}

		// Download Aim trainer contents
		aimTrainerScriptsOutput, err := download.StartLatestRepoContentsDownload(
			ghClient,
			errGroup,
			"Downloading AimTrainer Scripts",
			cacheDir,
			"scripts",
			"ColombianGuy",
			"r5_aimtrainer",
		)
		if err != nil {
			fileLogger.Error().Err(fmt.Errorf("error starting download of AimTrainer scripts: %v", err)).Msg("error")
			util.LogErrorWithDialog(fmt.Errorf("error starting download of AimTrainer scripts: %v", err))
			return
		}

		if err := errGroup.Wait(); err != nil {
			fileLogger.Error().Err(fmt.Errorf("error encountered while performing AimTrainer downloads: %v", err)).Msg("error")
			util.LogErrorWithDialog(fmt.Errorf("error encountered while performing AimTrainer downloads: %v", err))
			return
		}

		// Unzip AimTrainer deps into R5Folder
		err = util.UnzipFile(aimTrainerReleaseOutput, r5Folder, false, "Extracting AimTrainer Deps")
		if err != nil {
			fileLogger.Error().Err(fmt.Errorf("error unzipping AimTrainer deps: %v", err)).Msg("error")
			util.LogErrorWithDialog(fmt.Errorf("error unzipping AimTrainer deps: %v", err))
			return
		}

		//Unzip AimTrainer Scripts into platform/scripts
		err = util.UnzipFile(aimTrainerScriptsOutput, filepath.Join(r5Folder, "platform/scripts"), true, "Extracting AimTrainer Scripts")
		if err != nil {
			fileLogger.Error().Err(fmt.Errorf("error unzipping AimTrainer scripts: %v", err)).Msg("error")
			util.LogErrorWithDialog(fmt.Errorf("error unzipping AimTrainer scripts: %v", err))
			return
		}
	}

	_ = dialog.Raise("Success. Confirm to close terminal.")
	return
}
