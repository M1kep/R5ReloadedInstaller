package main

import (
	"R5ReloadedInstaller/pkg/download"
	"R5ReloadedInstaller/pkg/util"
	"R5ReloadedInstaller/pkg/validation"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
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
	if !validation.IsInR5Folder() {
		if !(len(os.Args) >= 2) {
			_ = dialog.Raise("Please move the R5RInstaller into your R5 Directory")
			return
		}

		pathFromArgs := os.Args[1]
		if validation.IsR5Folder(pathFromArgs) {
			r5Folder = pathFromArgs
		} else {
			_ = dialog.Raise("Please move the R5RInstaller into your R5 Directory or pass correct path via arguments")
			return
		}
	} else {
		var err error
		r5Folder, err = os.Getwd()
		if err != nil {
			log.Fatal(fmt.Errorf("error retrieving current directory while validating r5Path: %v", err))
		}
	}

	if validation.IsLauncherFileLocked(r5Folder) {
		_ = dialog.Raise("Please close the R5 Launcher before running")
		return
	}

	println(r5Folder)
	cacheDir := filepath.Join(r5Folder, "R5InstallerDirectory/cache")
	installerTempDir := filepath.Join(r5Folder, "R5InstallerDirectory/tempStore")

	// Initialize required file structure
	err := os.MkdirAll(cacheDir, 0777)
	if err != nil {
		log.Fatal(fmt.Errorf("error initializing installer directory %s: %v", cacheDir, err))
	}

	err = os.MkdirAll(installerTempDir, 0777)
	if err != nil {
		log.Fatal(fmt.Errorf("error initializing installer directory %s: %v", installerTempDir, err))
	}

	//options := map[string]string{
	//	"sdkAndScripts":  "SDK + Scripts",
	//	"aimTrainer":     "Aim Trainer",
	//	"ts_CleanScript": "(Troubleshooting) Clean Scripts - Deletes 'platform/scripts' prior to extracting",
	//}

	prompt := &survey.MultiSelect{
		Message: "Please select options",
		Options: []string{
			"SDK + Scripts",
			"Aim Trainer",
			"(Troubleshooting) Clean Scripts - Deletes 'platform/scripts' prior to extracting",
		},
		Default: []int{
			0,
		},
	}

	var selectedOptions []string
	err = survey.AskOne(prompt, &selectedOptions)
	if err != nil {
		log.Fatal(err)
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
		// Download Depot Zip
		wg.Add(1)
		sdkOutputPath, err := download.GetLatestRepoRelease(ghClient, &wg, "Downloading SDK", cacheDir, "sdk-depot", "depot.zip", "Mauler125", "r5sdk")
		if err != nil {
			panic(fmt.Errorf("error downloading sdk release: %v", err))
		}

		// Download
		repoOwner := "Mauler125"
		//sdkRepoName := "r5sdk"
		scriptsRepoName := "scripts_r5"
		//r5SdkReleases, _, err := ghClient.Repositories.ListReleases(context.Background(), repoOwner, sdkRepoName, &github.ListOptions{})
		//if err != nil {
		//	log.Fatal(fmt.Errorf("error listing releases for %s/%s: %v", repoOwner, sdkRepoName, err))
		//}
		//latestSdkRelease := r5SdkReleases[0]
		//downloadUrl := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/depot.zip", repoOwner, sdkRepoName, *latestSdkRelease.TagName)
		//sdkOutputPath := filepath.Join(cacheDir, fmt.Sprintf("sdk-depot_%s.zip", *latestSdkRelease.TagName))
		//println(downloadUrl)
		//wg.Add(1)
		//go func() {
		//	defer wg.Done()
		//	err := download.DownloadFile(sdkOutputPath, downloadUrl, "Downloading SDK Release")
		//	if err != nil {
		//		log.Fatal(fmt.Errorf("error downloading release(%s) for %s/%s: %v", *latestSdkRelease.TagName, repoOwner, sdkRepoName, err))
		//	}
		//}()

		// Download Repo Contents
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
			panic(fmt.Errorf("error downloading scripts: %v", err))
		}

		//scriptsRepo, _, err := ghClient.Repositories.Get(context.Background(), repoOwner, scriptsRepoName)
		//if err != nil {
		//	log.Fatal(fmt.Errorf("error listing branches for %s/%s: %v", repoOwner, scriptsRepoName, err))
		//}
		//
		//scriptsCommits, _, err := ghClient.Repositories.ListCommits(context.Background(), repoOwner, scriptsRepoName, &github.CommitsListOptions{
		//	ListOptions: github.ListOptions{
		//		PerPage: 1,
		//	},
		//})
		//if err != nil {
		//	log.Fatal(fmt.Errorf("error listing commits for %s/%s: %v", repoOwner, scriptsRepoName, err))
		//}
		//commitShortSHA := (*scriptsCommits[0].SHA)[0:7]
		//scriptsDownloadUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/zipball/%s", repoOwner, scriptsRepoName, *scriptsRepo.DefaultBranch)
		//scriptsOutputPath := filepath.Join(cacheDir, fmt.Sprintf("scripts-%s_%s.zip", *scriptsRepo.DefaultBranch, commitShortSHA))
		//wg.Add(1)
		//go func() {
		//	defer wg.Done()
		//	err := download.DownloadFile(scriptsOutputPath, scriptsDownloadUrl, "Downloading Scripts")
		//	if err != nil {
		//		log.Fatal(fmt.Errorf("error downloading repo contents from %s/%s for commit %s: %v", repoOwner, scriptsRepoName, commitShortSHA, err))
		//	}
		//}()
		wg.Wait()

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := util.UnzipFile(sdkOutputPath, r5Folder, false, "Extracting SDK")
			if err != nil {
				log.Fatal(fmt.Errorf("error unzipping file %s: %v", sdkOutputPath, err))
			}
		}()
		wg.Wait()

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := util.UnzipFile(scriptsRepoContentsOutput, filepath.Join(r5Folder, "platform/scripts"), true, "Extracting scripts")
			if err != nil {
				log.Fatal(fmt.Errorf("error unzipping file %s: %v", scriptsRepoContentsOutput, err))
			}
		}()
		wg.Wait()

		// Extract Depot Zip
		// Extract Repo contents
	}

	if util.Contains(selectedOptions, "Aim Trainer") {
		// Download Aim Trainer release
		// Download Aim trainer contents
	}
	//
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	err := download.DownloadFile("./depot.zip", "https://github.com/Mauler125/r5sdk/releases/download/v2.1_rc4/depot.zip", "Downloading SDK")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}()
	//println("IN between?")
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	err := download.DownloadFile("./S3_N1094.zip", "https://api.github.com/repos/Mauler125/scripts_r5/zipball/S3_N1094", "Downloading Scripts")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}()
	//
	wg.Wait()
}
