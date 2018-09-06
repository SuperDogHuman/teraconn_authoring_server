package lessonType

import "time"

type Lesson struct {
	ID           string    `json:"id"`
	AvatarID     string    `json:"-"`
	Avatar       Avatar    `json:"avatar" datastore:"-"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	DurationSec  float64   `json:"durationSec"`
	ThumbnailURL string    `json:"thumbnailURL" datastore:"-"`
	GraphicIDs   []string  `json:"graphicIDs" datastore:"-"`
	ViewCount    int64     `json:"viewCount"`
	Version      int64     `json:"version"`
	IsPublic     bool      `json:"isPublic"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
}

type Avatar struct {
	ID           string    `json:"id"`
	UserID       string    `json:"userID"`
	ThumbnailURL string    `json:"thumbnailURL" datastore:"-"`
	Name         string    `json:"name"`
	Version      int64     `json:"version"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
}

type LessonAuthor struct {
	ID       string `json:"id"`
	LessonID string `json:"lessonID"`
	UserID   string `json:"userID"`
	Role     string `json:"role"`
}

type Graphic struct {
	ID                string `json:"id"`
	GraphicCategoryID string `json:"graphicCategoryID"`
	UserID            string `json:"userID"`
	FileType          string `json:"fileType"`
	WidthPx           int    `json:"widthPx"`
	HeightPx          int    `json:"heightPx"`
	IsPublic          bool   `json:"isPublic"`
}

type LessonGraphic struct {
	ID         string    `json:"id"` /* same as ID of Lesson */
	GraphicIDs []string  `json:"graphicIDs"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}

type LessonVoiceText struct {
	FileID      string  `json:"fileID"`
	LessonID    string  `json:"lessonId"`
	DurationSec float64 `json:"durationSec"`
	Text        string  `json:"text"`
	IsTexted    bool    `json:"isTexted"`
	IsConverted bool    `json:"isConverted"`
}
