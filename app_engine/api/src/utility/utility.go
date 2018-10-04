package utility

import (
	"context"
	"regexp"

	"google.golang.org/appengine"
)

var xidRegexp = regexp.MustCompile("[0-9a-v]{20}")

// IsValidXIDs is check valid XID string format.
func IsValidXIDs(ids []string) bool {
	regExp := xidRegexp.Copy()
	for _, id := range ids {
		if !regExp.MatchString(id) {
			return false
		}
	}

	return true
}

// MaterialBucketName is return bucket name each environments.
func MaterialBucketName(ctx context.Context) string {
	bucketName := "teraconn_material"
	module := appengine.ModuleName(ctx)

	if module == "teraconnect-api" {
		return bucketName
	}

	return bucketName + "_development"
}

// RawVoiceBucketName is return bucket name each environments.
func RawVoiceBucketName(ctx context.Context) string {
	bucketName := "teraconn_raw_voice"
	module := appengine.ModuleName(ctx)

	if module == "teraconnect-api" {
		return bucketName
	}

	return bucketName + "_development"
}

// VoiceForTranscriptionBucketName is return bucket name each environments.
func VoiceForTranscriptionBucketName(ctx context.Context) string {
	bucketName := "teraconn_voice_for_transcription"
	module := appengine.ModuleName(ctx)

	if module == "teraconnect-api" {
		return bucketName
	}

	return bucketName + "_development"
}
