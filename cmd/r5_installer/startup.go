package main

import (
	"R5ReloadedInstaller/pkg/validation"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v47/github"
	"github.com/tawesoft/golib/v2/dialog"
	"golang.org/x/mod/semver"
	"os"
	"path/filepath"
	"time"
)

func getValidatedR5Folder() (validatedFolder string, err error) {
	isRunningInR5Folder, err := validation.IsRunningInR5Folder()
	if err != nil {
		return "", err
	}

	if isRunningInR5Folder {
		validatedFolder, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("error retrieving current directory while validating r5Path: %v", err)
		}

		return validatedFolder, nil
	}

	// Check CLI argument
	if !(len(os.Args) >= 2) {
		_ = dialog.Raise("Please move the R5RInstaller into your R5 Directory")
		return "", fmt.Errorf("not running from r5Folder and no argument provided")
	}

	pathFromArgs := os.Args[1]
	isPathR5Folder, err := validation.IsR5Folder(pathFromArgs)
	if err != nil {
		return "", err
	}
	if isPathR5Folder {
		validatedFolder = pathFromArgs
	} else {
		_ = dialog.Raise("Please move the R5RInstaller into your R5 Directory or pass correct path via arguments")
		return "", fmt.Errorf("not running from r5Folder and provided path argument is invalid")
	}

	return validatedFolder, nil
}

func initializeDirectories(r5Folder string) (cacheDir string, err error) {
	cacheDir = filepath.Join(r5Folder, "R5InstallerDirectory/cache")

	err = os.MkdirAll(cacheDir, 0777)
	if err != nil {
		err = fmt.Errorf("error initializing installer directory %s: %v", cacheDir, err)
		return
	}

	return
}

type UpdateCheckDetails struct {
	LastUpdateCheck         time.Time
	LastRetrievedReleaseTag string
}

func checkForUpdate(ghClient *github.Client, cacheDir string, currentVersion string) (shouldExit bool, message string, err error) {
	repoOwner := "M1kep"

	repoName := "R5ReloadedInstaller"
	updateCheckDetailsFromDisk, err := loadUpdateCheckDetails(cacheDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return true, "Error loading update details from disk", fmt.Errorf("error encountered loading update check details: %v", err)
		}
	}

	useUpdateDetailsCache := false
	// Update check should not happen within 10 minutes of each other
	nextUpdateCheckAt := updateCheckDetailsFromDisk.LastUpdateCheck.Add(time.Minute * 10)
	if time.Now().Before(nextUpdateCheckAt) {
		timeTillNextCheck := nextUpdateCheckAt.Sub(time.Now())
		if updateCheckDetailsFromDisk.LastRetrievedReleaseTag != "" {
			useUpdateDetailsCache = true
		} else {
			// If we don't have a cached tag, and the last update check was within 10 minutes, don't continue.
			return false, fmt.Sprintf("INFO: Last update check may have failed, waiting %s to check again.", timeTillNextCheck), nil
		}
	}

	if !semver.IsValid(currentVersion) {
		return true, "Current version is invalid", fmt.Errorf("invalid current version provided '%s'", currentVersion)
	}

	newUpdateCheckDetails := UpdateCheckDetails{}
	var latestVersionTag string
	if useUpdateDetailsCache {
		latestVersionTag = updateCheckDetailsFromDisk.LastRetrievedReleaseTag
		newUpdateCheckDetails.LastUpdateCheck = updateCheckDetailsFromDisk.LastUpdateCheck

		// If a new version was downloaded within the udpatecheck delay window(10 minutes)
		// Then we should use and persist the current version
		if semver.Compare(currentVersion, latestVersionTag) > 0 {
			latestVersionTag = currentVersion
		}
	} else {
		newUpdateCheckDetails.LastUpdateCheck = time.Now()
		repoReleases, _, err := ghClient.Repositories.ListReleases(context.Background(), repoOwner, repoName, &github.ListOptions{})
		if err != nil {
			saveDetailsErr := saveUpdateDetails(cacheDir, newUpdateCheckDetails)
			if saveDetailsErr != nil {
				return false, "", err
			}

			return false, "", fmt.Errorf("error listing releases for %s/%s: %v", repoOwner, repoName, err)
		}

		latestVersionTag = *(repoReleases[0].TagName)
	}

	if !semver.IsValid(latestVersionTag) {
		err := saveUpdateDetails(cacheDir, newUpdateCheckDetails)
		if err != nil {
			return false, "", err
		}

		return false, "", fmt.Errorf("invalid version from GitHub release '%s'", latestVersionTag)
	}

	newUpdateCheckDetails.LastRetrievedReleaseTag = latestVersionTag
	err = saveUpdateDetails(cacheDir, newUpdateCheckDetails)
	if err != nil {
		return false, "", err
	}

	if semver.Compare(currentVersion, latestVersionTag) < 0 {
		if semver.Major(latestVersionTag) > semver.Major(currentVersion) {
			return true, "New major version available. Browser will open to the following link after closing: https://github.com/M1kep/R5ReloadedInstaller/releases/latest", nil
		}

		return false, "New minor update is available. Consider downloading the latest release from https://github.com/M1kep/R5ReloadedInstaller/releases/latest", nil
	}
	return false, "", nil
}

func saveUpdateDetails(cacheDir string, details UpdateCheckDetails) error {
	jsonOut, err := json.Marshal(details)
	if err != nil {
		return fmt.Errorf("error marshalling details: %v", err)
	}

	err = os.WriteFile(filepath.Join(cacheDir, "updateCheckDetails.json"), jsonOut, 0777)
	if err != nil {
		return fmt.Errorf("error writing update details to disk: %v", err)
	}

	return nil
}

func loadUpdateCheckDetails(cacheDir string) (UpdateCheckDetails, error) {
	fileBytes, err := os.ReadFile(filepath.Join(cacheDir, "updateCheckDetails.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return UpdateCheckDetails{}, err
		} else {
			return UpdateCheckDetails{}, fmt.Errorf("error reading update details from disk: %v", err)
		}
	}

	details := UpdateCheckDetails{}
	err = json.Unmarshal(fileBytes, &details)
	if err != nil {
		return UpdateCheckDetails{}, fmt.Errorf("error unmarshalling update details: %v", err)
	}

	return details, nil
}
