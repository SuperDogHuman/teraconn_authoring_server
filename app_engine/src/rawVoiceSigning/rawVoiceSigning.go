package rawVoiceSigning

import (
	"cloudHelper"
	"net/http"

	"github.com/labstack/echo"
	"github.com/rs/xid"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// Get is get signing of raw voice files function.
func Get(c echo.Context) error {
	lessonID := c.QueryParam("lesson_id")
	if lessonID == "" {
		return c.JSON(http.StatusBadRequest, "lesson_id params was not found.")
	}

	ctx := appengine.NewContext(c.Request())
	bucketName := "teraconn_raw_voice"
	fileID := xid.New().String()
	fileName := lessonID + "-" + fileID + ".wav"
	contentType := "audio/wav"

	if err := cloudHelper.CreateObjectToGCS(ctx, bucketName, fileName, contentType, nil); err != nil {
		log.Errorf(ctx, err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if signedURL, err := cloudHelper.GetGCSSignedURL(ctx, bucketName, fileID, fileName, "PUT", contentType); err != nil {
		log.Errorf(ctx, err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		return c.JSON(http.StatusOK, rawWavSign{FileID: fileID, SignedURL: signedURL})
	}
}

type rawWavSign struct {
	FileID    string `json:"file_id"`
	SignedURL string `json:"signed_url"`
}
