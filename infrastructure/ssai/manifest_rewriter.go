package ssai

/**
 * SSAI Manifest Rewriter
 * Issue #26: SSAI (Server-Side Ad Insertion) Setup
 * 
 * Rewrites HLS/DASH manifests to include ad breaks
 */

import (
	"strings"
	"fmt"
)

// AdBreak represents an ad break to insert
type AdBreak struct {
	StartTime   float64 // seconds from start
	Duration    float64 // seconds
	AdSegments  []string // URLs of ad segments
	Type        string // "pre-roll", "mid-roll", "post-roll"
}

// RewriteHLSManifest inserts ad breaks into HLS manifest
func RewriteHLSManifest(manifest string, adBreaks []AdBreak) string {
	lines := strings.Split(manifest, "\n")
	var result []string
	
	result = append(result, "#EXTM3U")
	result = append(result, "#EXT-X-VERSION:3")
	
	// Insert pre-roll ads
	for _, adBreak := range adBreaks {
		if adBreak.Type == "pre-roll" {
			// Insert ad break before content
			result = append(result, fmt.Sprintf("#EXT-X-CUE-OUT:%.1f", adBreak.Duration))
			for _, segment := range adBreak.AdSegments {
				result = append(result, "#EXTINF:6.0,")
				result = append(result, segment)
			}
			result = append(result, "#EXT-X-CUE-IN")
		}
	}
	
	// Insert content segments
	for i, line := range lines {
		if strings.HasPrefix(line, "#EXTINF") || strings.HasPrefix(line, "http") {
			// Check if we need to insert mid-roll ad
			currentTime := float64(i * 6) // Approximate time
			for _, adBreak := range adBreaks {
				if adBreak.Type == "mid-roll" && 
				   currentTime >= adBreak.StartTime && 
				   currentTime < adBreak.StartTime+adBreak.Duration {
					// Insert mid-roll ad break
					result = append(result, fmt.Sprintf("#EXT-X-CUE-OUT:%.1f", adBreak.Duration))
					for _, segment := range adBreak.AdSegments {
						result = append(result, "#EXTINF:6.0,")
						result = append(result, segment)
					}
					result = append(result, "#EXT-X-CUE-IN")
				}
			}
			result = append(result, line)
		} else if strings.HasPrefix(line, "#EXT") {
			result = append(result, line)
		}
	}
	
	// Insert post-roll ads
	for _, adBreak := range adBreaks {
		if adBreak.Type == "post-roll" {
			result = append(result, fmt.Sprintf("#EXT-X-CUE-OUT:%.1f", adBreak.Duration))
			for _, segment := range adBreak.AdSegments {
				result = append(result, "#EXTINF:6.0,")
				result = append(result, segment)
			}
			result = append(result, "#EXT-X-CUE-IN")
		}
	}
	
	result = append(result, "#EXT-X-ENDLIST")
	return strings.Join(result, "\n")
}

// RewriteDASHManifest inserts ad breaks into DASH manifest
func RewriteDASHManifest(manifest string, adBreaks []AdBreak) string {
	// TODO: Implement DASH manifest rewriting
	// Insert ad periods between content periods
	return manifest
}

