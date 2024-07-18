package face

import (
	"fmt"
	"github.com/photoprism/photoprism/internal/thumb/crop"
	"os"
	"strconv"
)

func getEnvOrDefault(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func getEnvOrDefaultInt(k string, def int) int {
	if v, err := strconv.ParseInt(getEnvOrDefault(k, string(def)), 10, 32); err == nil {
		return int(v)
	}
	return def
}

func getEnvOrDefaultF64(k string, def float64) float64 {
	if v, err := strconv.ParseFloat(getEnvOrDefault(k, fmt.Sprintf("%v", def)), 64); err == nil {
		return v
	}
	return def
}

var CropSize = crop.Sizes[crop.Tile160]                                             // Face image crop size for FaceNet.
var OverlapThreshold = getEnvOrDefaultInt("PHOTOPRISM_FACE_OVERLAP", 42)            // Face area overlap threshold in percent.
var OverlapThresholdFloor = OverlapThreshold - 1                                    // Reduced overlap area to avoid rounding inconsistencies.
var ScoreThreshold = getEnvOrDefaultF64("PHOTOPRISM_FACE_SCORE", 9.0)               // Min face score.
var ClusterScoreThreshold = getEnvOrDefaultInt("PHOTOPRISM_FACE_CLUSTER_SCORE", 15) // Min score for faces forming a cluster.
var SizeThreshold = getEnvOrDefaultInt("PHOTOPRISM_FACE_SIZE", 50)                  // Min face size in pixels.
var ClusterSizeThreshold = getEnvOrDefaultInt("PHOTOPRISM_FACE_CLUSTER_SIZE", 50)   // Min size for faces forming a cluster in pixels.
var ClusterDist = getEnvOrDefaultF64("PHOTOPRISM_FACE_CLUSTER_DIST", 0.64)          // Similarity distance threshold of faces forming a cluster core.
var MatchDist = getEnvOrDefaultF64("PHOTOPRISM_FACE_MATCH_DIST", 0.46)              // Dist offset threshold for matching new faces with clusters.
var ClusterCore = getEnvOrDefaultInt("PHOTOPRISM_FACE_CLUSTER_CORE", 4)             // Min number of faces forming a cluster core.
var SampleThreshold = 2 * ClusterCore                                               // Threshold for automatic clustering to start.

// QualityThreshold returns the scale adjusted quality score threshold.
func QualityThreshold(scale int) (score float32) {
	score = float32(ScoreThreshold)

	// Smaller faces require higher quality.
	switch {
	case scale < 26:
		score += 26.0
	case scale < 32:
		score += 16.0
	case scale < 40:
		score += 11.0
	case scale < 50:
		score += 9.0
	case scale < 80:
		score += 6.0
	case scale < 110:
		score += 2.0
	}

	return score
}
