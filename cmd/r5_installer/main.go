package main

import (
	"R5ReloadedInstaller/internal/download"
	"R5ReloadedInstaller/pkg/util"
	"R5ReloadedInstaller/pkg/validation"
	"fmt"
	"github.com/google/go-github/v47/github"
	"github.com/gosuri/uiprogress"
	"github.com/tawesoft/golib/v2/dialog"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	var r5Folder string
	ghClient := github.NewClient(nil)

	r5Folder, err := getValidatedR5Folder()
	if err != nil {
		util.ExitOnError(err)
	}

	if validation.IsLauncherFileLocked(r5Folder) {
		_ = dialog.Raise("Please close the R5 Launcher before running.")
		return
	}

	cacheDir, err := initializeDirectories(r5Folder)
	if err != nil {
		util.ExitOnError(err)
	}

	selectedOptions, err := gatherRunOptions([]string{
		"SDK + Scripts",
		"Aim Trainer",
		"(Troubleshooting) Clean Scripts - Deletes 'platform/scripts' prior to extracting",
	})
	if err != nil {
		util.ExitOnError(fmt.Errorf("error gathering run options"))
	}

	uiprogress.Start()
	var wg sync.WaitGroup

	// PHASE 1
	if util.Contains(selectedOptions, "(Troubleshooting) Clean Scripts - Deletes 'platform/scripts' prior to extracting") {
		err := os.RemoveAll(filepath.Join(r5Folder, "platform/scripts"))
		if err != nil {
			log.Fatal(fmt.Errorf("error removing 'platform/scripts' folder: %v", err))
		}
	}

	// Phase 2
	if util.Contains(selectedOptions, "SDK + Scripts") {
		// Download SDK Release
		wg.Add(1)
		sdkOutputPath, err := download.GetLatestRepoRelease(
			ghClient,
			&wg,
			"Downloading SDK",
			cacheDir,
			"sdk-depot",
			"depot.zip",
			"Mauler125",
			"r5sdk",
		)
		if err != nil {
			util.ExitOnError(fmt.Errorf("error downloading sdk release: %v", err))
		}

		// Download scripts_r5
		wg.Add(1)
		scriptsRepoContentsOutput, err := download.GetLatestRepoContents(
			ghClient,
			&wg,
			"Downloading Scripts",
			cacheDir,
			"scripts",
			"Mauler125",
			"scripts_r5",
		)
		if err != nil {
			util.ExitOnError(fmt.Errorf("error downloading scripts: %v", err))
		}
		wg.Wait()

		// Unzip SDK into R5Folder
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := util.UnzipFile(sdkOutputPath, r5Folder, false, "Extracting SDK")
			if err != nil {
				util.ExitOnError(fmt.Errorf("error unzipping file %s: %v", sdkOutputPath, err))
			}
		}()
		wg.Wait()

		// Unzip Scripts into platform/scripts
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := util.UnzipFile(scriptsRepoContentsOutput, filepath.Join(r5Folder, "platform/scripts"), true, "Extracting scripts")
			if err != nil {
				util.ExitOnError(fmt.Errorf("error unzipping file %s: %v", scriptsRepoContentsOutput, err))
			}
		}()
		wg.Wait()
	}

	if util.Contains(selectedOptions, "Aim Trainer") {
		// Download Aim Trainer release
		wg.Add(1)
		aimTrainerReleaseOutput, err := download.GetLatestRepoRelease(
			ghClient,
			&wg,
			"Downloading Aim Trainer",
			cacheDir,
			"aimtrainer-deps",
			"AimTrainerRequiredFiles.zip",
			"ColombianGuy",
			"r5_aimtrainer",
		)
		if err != nil {
			util.ExitOnError(fmt.Errorf("error downloading aimtrainer release: %v", err))
		}

		// Download Aim trainer contents
		wg.Add(1)
		aimTrainerScriptsOutput, err := download.GetLatestRepoContents(
			ghClient,
			&wg,
			"Downloading AimTrainer Scripts",
			cacheDir,
			"scripts",
			"ColombianGuy",
			"r5_aimtrainer",
		)
		if err != nil {
			util.ExitOnError(fmt.Errorf("error downloading AimTrainer scripts: %v", err))
		}
		wg.Wait()

		// Unzip AimTrainer deps into R5Folder
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := util.UnzipFile(aimTrainerReleaseOutput, r5Folder, false, "Extracting AimTrainer Deps")
			if err != nil {
				util.ExitOnError(fmt.Errorf("error unzipping file %s: %v", aimTrainerReleaseOutput, err))
			}
		}()
		wg.Wait()

		//Unzip AimTrainer Scripts into platform/scripts
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := util.UnzipFile(aimTrainerScriptsOutput, filepath.Join(r5Folder, "platform/scripts"), true, "Extracting AimTrainer Scripts")
			if err != nil {
				util.ExitOnError(fmt.Errorf("error unzipping file %s: %v", aimTrainerScriptsOutput, err))
			}
		}()
		wg.Wait()
	}

	util.ExitWithAlertDialog(false)
}
