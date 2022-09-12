package download

import (
	"context"
	"fmt"
	"github.com/google/go-github/v47/github"
	"log"
	"path/filepath"
	"sync"
)

func GetLatestRepoRelease(ghClient *github.Client, wg *sync.WaitGroup, progressMessage string, cacheDirectory string, cacheName string, releaseFileName string, repoOwner string, repoName string) (outputPath string, err error) {
	repoReleases, _, err := ghClient.Repositories.ListReleases(context.Background(), repoOwner, repoName, &github.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("error listing releases for %s/%s: %v", repoOwner, repoName, err)
	}

	latestRelease := repoReleases[0]
	downloadUrl := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s", repoOwner, repoName, *latestRelease.TagName, releaseFileName)
	sdkOutputPath := filepath.Join(cacheDirectory, fmt.Sprintf("%s_%s.zip", cacheName, *latestRelease.TagName))

	go func() {
		defer wg.Done()
		err := DownloadFile(sdkOutputPath, downloadUrl, progressMessage)
		if err != nil {
			log.Fatal(fmt.Errorf("error downloading release(%s) for %s/%s: %v", *latestRelease.TagName, repoOwner, repoName, err))
		}
		print(fmt.Sprintf("Finished downloading: %s", downloadUrl))
	}()

	return sdkOutputPath, nil
}

func GetLatestRepoContents(ghClient *github.Client, wg *sync.WaitGroup, progressMessage string, cacheDirectory string, cacheName string, repoOwner string, repoName string) (outputPath string, err error) {
	ghRepo, _, err := ghClient.Repositories.Get(context.Background(), repoOwner, repoName)
	if err != nil {
		return "", fmt.Errorf("error retrieving repo info for %s/%s: %v", repoOwner, repoName, err)
	}

	repoCommits, _, err := ghClient.Repositories.ListCommits(context.Background(), repoOwner, repoName, &github.CommitsListOptions{
		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	})
	if err != nil {
		return "", fmt.Errorf("error retrieving commits for %s/%s: %v", repoOwner, repoName, err)
	}
	latestCommitShortSHA := (*repoCommits[0].SHA)[0:7]
	downloadUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/zipball/%s", repoOwner, repoName, *ghRepo.DefaultBranch)
	outputPath = filepath.Join(cacheDirectory, fmt.Sprintf("%s-%s_%s.zip", cacheName, *ghRepo.DefaultBranch, latestCommitShortSHA))

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := DownloadFile(outputPath, downloadUrl, progressMessage)
		if err != nil {
			log.Fatal(fmt.Errorf("error downloading repo contents from %s/%s for commit %s: %v", repoOwner, repoName, latestCommitShortSHA, err))
		}
	}()

	return outputPath, nil
}
