package voiceText

import (
	"context"
	"lessonType"
	"net/http"

	"github.com/labstack/echo"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

// Gets is get texts from voice function.
func Gets(c echo.Context) error {
	ctx := appengine.NewContext(c.Request())
	id := c.Param("lesson_id")

	if voiceTexts, err := fetchVoiceTextsFromGCD(ctx, id); err != nil {
		// TODO return http status each error type
		log.Errorf(ctx, err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		return c.JSON(http.StatusOK, voiceTexts)
	}
}

func fetchVoiceTextsFromGCD(ctx context.Context, lessonID string) ([]lessonType.VoiceText, error) {
	query := datastore.NewQuery("VoiceText").Filter("LessonID =", lessonID).Order("FileID")

	var voiceTexts []lessonType.VoiceText
	if _, err := query.GetAll(ctx, &voiceTexts); err != nil {
		return voiceTexts, err
	}

	return voiceTexts, nil
}
