package utility

import (
	"context"
	"regexp"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
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
	module := appengine.ModuleName(ctx)
	log.Infof(ctx, module)

	if module == "teraconnect-api" {
		return "teraconn_material"
	}

	return "teraconn_material_development"
}

// RawVoiceBucketName is return bucket name each environments.
func RawVoiceBucketName(ctx context.Context) string {
	module := appengine.ModuleName(ctx)
	log.Infof(ctx, module)

	if module == "teraconnect-api" {
		return "teraconn_raw_voice"
	}

	return "teraconn_raw_voice_development"
}
