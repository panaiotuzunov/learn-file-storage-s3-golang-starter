package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func getAssetPath(randomString string, mediaType string) string {
	ext := mediaTypeToExt(mediaType)
	return fmt.Sprintf("%s%s", randomString, ext)
}

func (cfg apiConfig) getAssetDiskPath(assetPath string) string {
	return filepath.Join(cfg.assetsRoot, assetPath)
}

func (cfg apiConfig) getAssetURL(assetPath string) string {
	return fmt.Sprintf("http://localhost:%s/assets/%s", cfg.port, assetPath)
}

func mediaTypeToExt(mediaType string) string {
	parts := strings.Split(mediaType, "/")
	if len(parts) != 2 {
		return ".bin"
	}
	return "." + parts[1]
}

func getVideoAspectRatio(filePath string) (string, error) {
	command := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	result := bytes.Buffer{}
	command.Stdout = &result
	if err := command.Run(); err != nil {
		return "", err
	}

	var data struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}
	if err := json.Unmarshal(result.Bytes(), &data); err != nil {
		return "", err
	}

	width := data.Streams[0].Width
	height := data.Streams[0].Height

	currentRatio := float64(width) / float64(height)
	targetRatioPortrait := float64(9) / float64(16)
	targetRatioLandscape := float64(16) / float64(9)
	tolerance := 0.05

	switch {
	case math.Abs(currentRatio-targetRatioLandscape) < tolerance:
		return "landscape", nil
	case math.Abs(currentRatio-targetRatioPortrait) < tolerance:
		return "portrait", nil
	default:
		return "other", nil
	}
}
