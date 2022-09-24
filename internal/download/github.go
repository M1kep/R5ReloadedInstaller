package download

import (
	"R5ReloadedInstaller/pkg/download"
	"context"
	"fmt"
	"github.com/google/go-github/v47/github"
	"golang.org/x/sync/errgroup"
	"path/filepath"
)

func StartLatestRepoReleaseDownload(ghClient *github.Client, eg *errgroup.Group, progressMessage string, cacheDirectory string, cacheName string, releaseFileName string, repoOwner string, repoName string, includePreReleases bool) (outputPath string, err error) {
	repoReleases, _, err := ghClient.Repositories.ListReleases(context.Background(), repoOwner, repoName, &github.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("error listing releases for %s/%s: %v", repoOwner, repoName, err)
	}

	releaseToDownload := repoReleases[0]
	if !includePreReleases {
		for _, repo := range repoReleases {
			if !*repo.Prerelease {
				releaseToDownload = repo
				break
			}
		}
	}

	downloadUrl := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s", repoOwner, repoName, *releaseToDownload.TagName, releaseFileName)
	sdkOutputPath := filepath.Join(cacheDirectory, fmt.Sprintf("%s_%s.zip", cacheName, *releaseToDownload.TagName))

	eg.Go(func() error {
		err := download.DownloadFile(sdkOutputPath, downloadUrl, progressMessage)
		if err != nil {
			return fmt.Errorf("error downloading release(%s) for %s/%s: %v", *releaseToDownload.TagName, repoOwner, repoName, err)
		}
		return nil
	})

	return sdkOutputPath, nil
}

func StartLatestRepoContentsDownload(ghClient *github.Client, eg *errgroup.Group, progressMessage string, cacheDirectory string, cacheName string, repoOwner string, repoName string) (outputPath string, err error) {
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

	eg.Go(func() error {
		err := download.DownloadFile(outputPath, downloadUrl, progressMessage)
		if err != nil {
			return fmt.Errorf("error downloading repo contents from %s/%s for commit %s: %v", repoOwner, repoName, latestCommitShortSHA, err)
		}
		return nil
	})

	return outputPath, nil
}
