package lessonType

import "time"

type Lesson struct {
	ID           string    `json:"id" datastore:"-"`
	AvatarID     string    `json:"avatarID"`
	Avatar       Avatar    `json:"avatar" datastore:"-"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	DurationSec  float64   `json:"durationSec"`
	ThumbnailURL string    `json:"thumbnailURL" datastore:"-"`
	GraphicIDs   []string  `json:"graphicIDs"`
	Graphics     []Graphic `json:"graphics" datastore:"-"`
	ViewCount    int64     `json:"viewCount"`
	Version      int64     `json:"version"`
	IsPacked     bool      `json:"isPacked"`
	IsPublic     bool      `json:"isPublic"`
	UserID       string    `json:"userID"`
	ShouldDelete bool      `json:"-"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
}

type Avatar struct {
	ID           string    `json:"id" datastore:"-"`
	UserID       string    `json:"userID"`
	URL          string    `json:"url" datastore:"-"`
	ThumbnailURL string    `json:"thumbnailURL" datastore:"-"`
	Name         string    `json:"name"`
	Version      int64     `json:"version"`
	IsPublic     bool      `json:"isPublic"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
}

type LessonAuthor struct {
	ID       string    `json:"id" datastore:"-"`
	LessonID string    `json:"lessonID"`
	UserID   string    `json:"userID"`
	Role     string    `json:"role"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

type Graphic struct {
	ID                string    `json:"id" datastore:"-"`
	URL               string    `json:"url" datastore:"-"`
	ThumbnailURL      string    `json:"thumbnailURL" datastore:"-"`
	GraphicCategoryID string    `json:"graphicCategoryID"`
	UserID            string    `json:"userID"`
	FileType          string    `json:"fileType"`
	IsPublic          bool      `json:"isPublic"`
	Created           time.Time `json:"created"`
	Updated           time.Time `json:"updated"`
}

type LessonVoiceText struct {
	FileID      string  `json:"fileID"`
	LessonID    string  `json:"lessonID"`
	DurationSec float64 `json:"durationSec"`
	Text        string  `json:"text"`
	IsTexted    bool    `json:"isTexted"`
	IsConverted bool    `json:"isConverted"`
}
