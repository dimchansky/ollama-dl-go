package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

const (
	numRetries = 10
)

var mediaTypeToFileTemplate = map[string]string{
	"application/vnd.ollama.image.license":  "license-%s.txt",
	"application/vnd.ollama.image.model":    "model-%s.gguf",
	"application/vnd.ollama.image.params":   "params-%s.json",
	"application/vnd.ollama.image.system":   "system-%s.txt",
	"application/vnd.ollama.image.template": "template-%s.txt",
}

type Layer struct {
	MediaType string `json:"mediaType"`
	Digest    string `json:"digest"`
	Size      int64  `json:"size"`
}

type Manifest struct {
	MediaType string  `json:"mediaType"`
	Layers    []Layer `json:"layers"`
}

type DownloadJob struct {
	Layer    Layer
	DestPath string
	BlobURL  string
	Size     int64
}

func getShortHash(layer Layer) (string, error) {
	if !strings.HasPrefix(layer.Digest, "sha256:") {
		return "", fmt.Errorf("unexpected digest: %s", layer.Digest)
	}
	return layer.Digest[7:19], nil
}

func downloadBlob(client *http.Client, job DownloadJob, wg *sync.WaitGroup) error {
	defer wg.Done()

	for attempt := 1; attempt <= numRetries; attempt++ {
		tempPath := job.DestPath + ".tmp"

		// Ensure the directory exists
		if err := os.MkdirAll(filepath.Dir(tempPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}

		outFile, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		defer outFile.Close()

		// Check for partial download
		startOffset, _ := outFile.Seek(0, io.SeekEnd)
		req, err := http.NewRequest("GET", job.BlobURL, nil)
		if err != nil {
			return err
		}

		if startOffset > 0 {
			req.Header.Set("Range", fmt.Sprintf("bytes=%d-", startOffset))
		}

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		bar := progressbar.DefaultBytes(job.Size, job.DestPath)
		bar.Set64(startOffset)

		_, err = io.Copy(io.MultiWriter(outFile, bar), resp.Body)
		if err != nil {
			continue
		}

		// Rename the temporary file to the final destination
		if err := os.Rename(tempPath, job.DestPath); err != nil {
			return err
		}
		return nil
	}

	return errors.New("maximum retries reached")
}

func getDownloadJobs(client *http.Client, registry, destDir, name, version string) ([]DownloadJob, error) {
	manifestURL := fmt.Sprintf("%s/v2/%s/manifests/%s", registry, name, version)
	resp, err := client.Get(manifestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get manifest: %d", resp.StatusCode)
	}

	var manifest Manifest
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, err
	}

	if manifest.MediaType != "application/vnd.docker.distribution.manifest.v2+json" {
		return nil, fmt.Errorf("unexpected media type for manifest: %s", manifest.MediaType)
	}

	var jobs []DownloadJob
	for _, layer := range manifest.Layers {
		fileTemplate, ok := mediaTypeToFileTemplate[layer.MediaType]
		if !ok {
			continue
		}

		shortHash, err := getShortHash(layer)
		if err != nil {
			return nil, err
		}

		filename := fmt.Sprintf(fileTemplate, shortHash)
		destPath := filepath.Join(destDir, filename)
		blobURL := fmt.Sprintf("%s/v2/%s/blobs/%s", registry, name, layer.Digest)

		jobs = append(jobs, DownloadJob{
			Layer:    layer,
			DestPath: destPath,
			BlobURL:  blobURL,
			Size:     layer.Size,
		})
	}

	return jobs, nil
}

func main() {
	registry := flag.String("registry", "https://registry.ollama.ai/", "Registry URL")
	destDir := flag.String("d", "", "Destination directory")

	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Usage: ollama-dl <name>")
		os.Exit(1)
	}

	name := flag.Arg(0)
	if !strings.Contains(name, "/") {
		name = "library/" + name
	}

	// Check for version and append ":latest" if not specified
	if !strings.Contains(name, ":") {
		name += ":latest"
	}

	// Construct the destination directory name after handling the version
	if *destDir == "" {
		*destDir = strings.ReplaceAll(strings.ReplaceAll(name, "/", "-"), ":", "-")
	}

	nameParts := strings.Split(name, ":")
	name, version := nameParts[0], nameParts[1]

	client := &http.Client{Timeout: 30 * time.Second}
	jobs, err := getDownloadJobs(client, *registry, *destDir, name, version)
	if err != nil {
		fmt.Println("Error getting download jobs:", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	for _, job := range jobs {
		if _, err := os.Stat(job.DestPath); err == nil {
			fmt.Println("Already have", job.DestPath)
			continue
		}
		wg.Add(1)
		go func(job DownloadJob) {
			if err := downloadBlob(client, job, &wg); err != nil {
				fmt.Println("Download error:", err)
			}
		}(job)
	}

	wg.Wait()
	fmt.Println("Download complete")
}
