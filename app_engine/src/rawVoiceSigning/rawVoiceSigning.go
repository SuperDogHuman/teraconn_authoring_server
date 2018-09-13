package rawVoiceSigning

import (
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"cloudHelper"
	"net/http"
	"utility"
)

// Get is get signing of raw voice files function.
func Get(c echo.Context) error {
	lessonID := c.QueryParam("lesson_id")
	if lessonID == "" {
		return c.JSON(http.StatusBadRequest, "lesson_id params was not found.")
	}

	ctx := appengine.NewContext(c.Request())

	ids := []string{lessonID}
	if !utility.IsValidXIDs(ids) {
		errMessage := "Invalid ID(s) error"
		log.Warningf(ctx, errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	bucketName := "teraconn_raw_voice"
	fileID := xid.New().String()
	fileName := lessonID + "-" + fileID + ".wav"
	contentType := "audio/wav"

	if err := cloudHelper.CreateObjectToGCS(ctx, bucketName, fileName, contentType, nil); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if signedURL, err := cloudHelper.GetGCSSignedURL(ctx, bucketName, fileName, "PUT", contentType); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		return c.JSON(http.StatusOK, rawWavSign{FileID: fileID, SignedURL: signedURL})
	}
}

type rawWavSign struct {
	FileID    string `json:"file_id"`
	SignedURL string `json:"signed_url"`
}
