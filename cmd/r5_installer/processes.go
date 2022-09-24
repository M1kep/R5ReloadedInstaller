package main

import (
	"R5ReloadedInstaller/internal/download"
	"R5ReloadedInstaller/pkg/util"
	"fmt"
	"github.com/google/go-github/v47/github"
	"golang.org/x/sync/errgroup"
	"path/filepath"
)

func ProcessSDK(ghClient *github.Client, errGroup *errgroup.Group, cacheDir string, r5Folder string, includePreReleases bool) error {
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
		includePreReleases,
	)
	if err != nil {
		return fmt.Errorf("error starting download of sdk release: %v", err)
	}

	if err := errGroup.Wait(); err != nil {
		return fmt.Errorf("error encountered while performing SDK download: %v", err)
	}

	// Unzip SDK into R5Folder
	err = util.UnzipFile(sdkOutputPath, r5Folder, false, "Extracting SDK")
	if err != nil {
		return fmt.Errorf("error unzipping sdk: %v", err)
	}

	return nil
}

func ProcessLatestR5Scripts(ghClient *github.Client, errGroup *errgroup.Group, cacheDir string, r5Folder string) error {
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
		return fmt.Errorf("error starting download of scripts: %v", err)
	}

	if err := errGroup.Wait(); err != nil {
		return fmt.Errorf("error encountered while performing r5_scripts download: %v", err)
	}

	// Unzip Scripts into platform/scripts
	err = util.UnzipFile(scriptsRepoContentsOutput, filepath.Join(r5Folder, "platform/scripts"), true, "Extracting scripts")
	if err != nil {
		return fmt.Errorf("error unzipping scripts: %v", err)
	}

	return nil
}

func ProcessAimTrainer(ghClient *github.Client, errGroup *errgroup.Group, cacheDir string, r5Folder string) error {
	aimTrainerReleaseOutput, err := download.StartLatestRepoReleaseDownload(
		ghClient,
		errGroup,
		"Downloading Aim Trainer",
		cacheDir,
		"aimtrainer-deps",
		"Flowstate.-.Required.Files.zip",
		"ColombianGuy",
		"r5_aimtrainer",
		false,
	)
	if err != nil {
		return fmt.Errorf("error starting download of aimtrainer release: %v", err)
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
		return fmt.Errorf("error starting download of AimTrainer scripts: %v", err)
	}

	if err := errGroup.Wait(); err != nil {
		return fmt.Errorf("error encountered while performing AimTrainer downloads: %v", err)
	}

	// Unzip AimTrainer deps into R5Folder
	err = util.UnzipFile(aimTrainerReleaseOutput, r5Folder, false, "Extracting AimTrainer Deps")
	if err != nil {
		return fmt.Errorf("error unzipping AimTrainer deps: %v", err)
	}

	//Unzip AimTrainer Scripts into platform/scripts
	err = util.UnzipFile(aimTrainerScriptsOutput, filepath.Join(r5Folder, "platform/scripts"), true, "Extracting AimTrainer Scripts")
	if err != nil {
		return fmt.Errorf("error unzipping AimTrainer scripts: %v", err)
	}

	return nil
}
