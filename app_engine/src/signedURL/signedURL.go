package signedURL

import (
	"cloudHelper"
	"net/http"

	"github.com/labstack/echo"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

const bucketName = "teraconn_material"

// Gets is get signed URLs of files.
func Gets(c echo.Context) error {
	params := new(GetsRequestParams)
	if err := c.Bind(params); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	ctx := appengine.NewContext(c.Request())
	urlLength := len(params.FileRequests)
	urls := make([]string, urlLength)
	for _, fileRequest := range params.FileRequests {
		// TODO check user permission
		// TODO check file exists

		filePath := fileRequest.Entity + "/" + fileRequest.ID + "." + fileRequest.Extension
		signedURL, err := cloudHelper.GetGCSSignedURL(ctx, bucketName, filePath, "GET", "")
		if err != nil {
			log.Errorf(ctx, err.Error())
		}
		urls = append(urls, signedURL)
	}

	return c.JSON(http.StatusOK, GetsResponses{SignedURLs: urls})
}

type GetsRequestParams struct {
	FileRequests []FileRequest `json:"file_requests"`
}

type FileRequest struct {
	ID        string `json:"id"`
	Entity    string `json:"entity"`
	Extension string `json:"extension"`
}

type GetsResponses struct {
	SignedURLs []string `json:"signed_urls"`
}
