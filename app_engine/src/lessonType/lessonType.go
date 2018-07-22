package lessonType

import "time"

type Lesson struct {
	ID           string    `json:"id" datastore:"-"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	DurationMSec int32     `json:"duration_msec"`
	ViewCount    int64     `json:"view_count"`
	ThumbnailURL string    `json:"thumbnail_url"`
	GraphicIDs   []string  `json:"graphic_ids"`
	Published    time.Time `json:"published"`
	Updated      time.Time `json:"updated"`
}

type LessonAuthor struct {
	ID       string `json:"id" datastore:"-"`
	LessonID string `json:"lesson_id"`
	UserID   string `json:"user_id"`
	Role     string `json:"role"`
}

type Graphic struct {
	ID     string   `json:"id" datastore:"-"`
	UserID string   `json:"user_id"`
	TeamID []string `json:"owner_ids"`
}

/* The following structs is for json.Unmarshall */

type LessonMaterial struct {
	DurationMSec uint32             `json:"duration_msec"`
	Timelines    []LessonTimeline   `json:"timelines"`
	Poses        []LessonAvatarPose `json:"poses"`
	Published    time.Time          `json:"published"`
	Updated      time.Time          `json:"updated"`
}

type LessonTimeline struct {
	TimeMSec uint32                    `json:"time_msec"`
	Text     LessonText                `json:"text"`
	Voice    LessonVoice               `json:"voice"`
	Graphic  []LessonGraphic           `json:"graphics"`
	SPAction LessonAvatarSpecialAction `json:"sp_action"`
}

type LessonText struct {
	DurationMSec uint32 `json:"duration_msec"`
	DelayMSec    uint32 `json:"delay_msec"`
	Text         string `json:"text"`
	Position     string `json:"position"`
	Style        string `json:"style"`
	Size         uint8  `json:"size"`
	Color        string `json:"color"`
}

type LessonVoice struct {
	FileID       string `json:"file_id"`
	DurationMSec uint32 `json:"duration_msec"`
	DelayMSec    uint32 `json:"delay_msec"`
}

type LessonGraphic struct {
	DelayMSec uint32 `json:"delay_msec"`
	GraphicID string `json:"graphic_id"`
	Action    string `json:"action"`
	WidthPx   uint16 `json:"width_px"`
	HeightPx  uint16 `json:"height_px"`
	Position  string `json:"position"`
}

type LessonAvatarSpecialAction struct {
	Action           string `json:"action"`
	FacialExpression string `json:"facial_expression"`
}

type LessonAvatarPose struct {
	TimeMSec   uint32                     `json:"time_msec"`
	LeftHand   LessonAvatarPositionVector `json:"left_hand"`
	RightHand  LessonAvatarPositionVector `json:"right_hand"`
	LeftElbow  LessonAvatarPositionVector `json:"left_elbow"`
	RightElbow LessonAvatarPositionVector `json:"right_elbow"`
	LookAt     LessonAvatarPositionVector `json:"look_at"`
	CoreBody   LessonAvatarPositionVector `json:"core_body"`
}

type LessonAvatarPositionVector struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}
