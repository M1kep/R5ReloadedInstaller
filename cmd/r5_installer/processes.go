package main

import (
	"R5ReloadedInstaller/internal/download"
	"R5ReloadedInstaller/pkg/util"
	"fmt"
	"github.com/google/go-github/v52/github"
	"golang.org/x/sync/errgroup"
	"path/filepath"
)

type ProcessManager struct {
	ghClient *github.Client
	errGroup *errgroup.Group
	cacheDir string
	r5Folder string
}

func (pm *ProcessManager) ProcessSDK(includePreReleases bool) error {
	// Download SDK Release
	sdkOutputPath, err := download.StartLatestRepoReleaseDownload(
		pm.ghClient,
		pm.errGroup,
		"Downloading SDK",
		pm.cacheDir,
		"sdk-depot",
		"depot.zip",
		"Mauler125",
		"r5sdk",
		includePreReleases,
	)
	if err != nil {
		return fmt.Errorf("error starting download of sdk release: %v", err)
	}

	if err := pm.errGroup.Wait(); err != nil {
		return fmt.Errorf("error encountered while performing SDK download: %v", err)
	}

	// Unzip SDK into R5Folder
	err = util.UnzipFile(sdkOutputPath, pm.r5Folder, false, "Extracting SDK")
	if err != nil {
		return fmt.Errorf("error unzipping sdk: %v", err)
	}

	return nil
}

func (pm *ProcessManager) ProcessLatestR5Scripts() error {
	// Download scripts_r5
	scriptsRepoContentsOutput, err := download.StartLatestRepoContentsDownload(
		pm.ghClient,
		pm.errGroup,
		"Downloading Scripts",
		pm.cacheDir,
		"scripts",
		"Mauler125",
		"scripts_r5",
	)
	if err != nil {
		return fmt.Errorf("error starting download of scripts: %v", err)
	}

	if err := pm.errGroup.Wait(); err != nil {
		return fmt.Errorf("error encountered while performing r5_scripts download: %v", err)
	}

	// Unzip Scripts into platform/scripts
	err = util.UnzipFile(scriptsRepoContentsOutput, filepath.Join(pm.r5Folder, "platform/scripts"), true, "Extracting scripts")
	if err != nil {
		return fmt.Errorf("error unzipping scripts: %v", err)
	}

	return nil
}

func (pm *ProcessManager) ProcessFlowstate() error {
	flowstateReleaseOutput, err := download.StartLatestRepoReleaseDownload(
		pm.ghClient,
		pm.errGroup,
		"Downloading FlowState Required Files",
		pm.cacheDir,
		"flowstate-deps",
		"Flowstate.-.Required.Files.zip",
		"ColombianGuy",
		"r5_flowstate",
		false,
	)
	if err != nil {
		return fmt.Errorf("error starting download of Flowstate release: %v", err)
	}

	// Download Aim trainer contents
	flowstateScriptsOutput, err := download.StartLatestRepoContentsDownload(
		pm.ghClient,
		pm.errGroup,
		"Downloading Latest Flowstate Scripts",
		pm.cacheDir,
		"scripts",
		"ColombianGuy",
		"r5_flowstate",
	)
	if err != nil {
		return fmt.Errorf("error starting download of Flowstate scripts: %v", err)
	}

	if err := pm.errGroup.Wait(); err != nil {
		return fmt.Errorf("error encountered while performing Flowstate downloads: %v", err)
	}

	// Unzip Flowstate deps into R5Folder
	err = util.UnzipFile(flowstateReleaseOutput, pm.r5Folder, false, "Extracting Flowstate Deps")
	if err != nil {
		return fmt.Errorf("error unzipping Flowstate deps: %v", err)
	}

	//Unzip Flowstate Scripts into platform/scripts
	err = util.UnzipFile(flowstateScriptsOutput, filepath.Join(pm.r5Folder, "platform/scripts"), true, "Extracting Flowstate Scripts")
	if err != nil {
		return fmt.Errorf("error unzipping Flowstate scripts: %v", err)
	}

	return nil
}
