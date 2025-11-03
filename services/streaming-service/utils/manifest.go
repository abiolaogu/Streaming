package utils

import (
	"fmt"
	"github.com/streamverse/streaming-service/models"
)

// GenerateHLSManifest generates an HLS manifest (.m3u8)
func GenerateHLSManifest(baseURL string, variants []models.Variant, subtitles []models.Subtitle, drmConfig *models.DRMConfig) string {
	manifest := "#EXTM3U\n"
	manifest += "#EXT-X-VERSION:6\n"
	manifest += "#EXT-X-INDEPENDENT-SEGMENTS\n"

	if drmConfig != nil {
		manifest += fmt.Sprintf("#EXT-X-KEY:METHOD=SAMPLE-AES,URI=\"%s\",KEYFORMAT=\"%s\",KEYFORMATVERSIONS=\"1\"\n",
			drmConfig.LicenseURL, drmConfig.Type)
	}

	// Add variants
	for _, variant := range variants {
		manifest += fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%s,CODECS=\"%s\"\n",
			variant.Bandwidth, variant.Resolution, variant.Codec)
		manifest += variant.URL + "\n"
	}

	// Add subtitles
	for _, subtitle := range subtitles {
		manifest += fmt.Sprintf("#EXT-X-MEDIA:TYPE=SUBTITLES,GROUP-ID=\"subs\",NAME=\"%s\",LANGUAGE=\"%s\",URI=\"%s\"\n",
			subtitle.Label, subtitle.Language, subtitle.URL)
	}

	return manifest
}

// GenerateDASHManifest generates a DASH manifest (.mpd)
func GenerateDASHManifest(baseURL string, variants []models.Variant, subtitles []models.Subtitle, drmConfig *models.DRMConfig) string {
	manifest := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"
	manifest += "<MPD xmlns=\"urn:mpeg:dash:schema:mpd:2011\" type=\"static\" mediaPresentationDuration=\"PT0H10M0S\" minBufferTime=\"PT1.5S\" profiles=\"urn:mpeg:dash:profile:isoff-on-demand:2011\">\n"
	manifest += "<Period>\n"
	manifest += "<AdaptationSet>\n"

	// Add variants
	for _, variant := range variants {
		manifest += fmt.Sprintf("<Representation id=\"%s\" bandwidth=\"%d\" width=\"1920\" height=\"1080\" codecs=\"%s\">\n",
			variant.Resolution, variant.Bandwidth, variant.Codec)
		manifest += fmt.Sprintf("<BaseURL>%s</BaseURL>\n", variant.URL)
		manifest += "</Representation>\n"
	}

	manifest += "</AdaptationSet>\n"
	manifest += "</Period>\n"
	manifest += "</MPD>\n"

	return manifest
}

