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

type LessonGraphic struct {
  ID     string `json:"id" datastore:"-"`
  UserID string `json:"user_id"`
}

/* The following structs is for json.Unmarshall */

type LessonMaterial struct {
  DurationMSec uint32                `json:"duration_msec"`
  Speeches     []LessonSpeech        `json:"speeches"`
  Actions      []LessonAvatarAction  `json:"actions"`
  Positions    []LessonAvatarPosition`json:"position"`
  PublishedAt  time.Time             `json:"published_at"`
  UpdatedAt    time.Time             `json:"updated_at"`
}

type LessonSpeech struct {
  SpeechedAtMSec uint32              `json:"speeched_at_msec"`
  DurationMSec   uint32              `json:"duration_msec"`
  Text           LessonSpeechText    `json:"text"`
  Voice          LessonSpeechVoice   `json:"voice"`
  Images         []LessonSpeechImage `json:"images"`
}

type LessonSpeechText struct {
  Text      string  `json:"text"`
  Position  string  `json:"position"`
  Style     string  `json:"style"`
  Size      uint8   `json:"size"`
  Color     string  `json:"color"`
  DelayMSec uint32  `json:"delay_msec"`
}

type LessonSpeechVoice struct {
  ID        string `json:"id"`
  DelayMSec uint16 `json:"delay_msec"`
}

type LessonSpeechImage struct {
  ID        string `json:"id"`
  DelayMSec uint32 `json:"delay_msec"`
  Action    string `json:"action"`
  WidthPx   uint16 `json:"width_px"`
  HeightPx  uint16 `json:"height_px"`
  Position  string `json:"position"`
}

type LessonAvatarAction struct {
  ActionedAtMSec uint32               `json:"actioned_at_msec"`
  Action         string               `json:"action"`
  Facial         string               `json:"facial"`
  Position       LessonAvatarPosition `json:"position"`
}

type LessonAvatarPosition struct {
  DurationMSec uint32                     `json:"duration_msec"`
  LeftHand     LessonAvatarPositionVector `json:"left_hand"`
  RightHand    LessonAvatarPositionVector `json:"right_hand"`
  LeftElbow    LessonAvatarPositionVector `json:"left_elbow"`
  RightElbow   LessonAvatarPositionVector `json:"right_elbow"`
  LookAt       LessonAvatarPositionVector `json:"look_at"`
  CoreBody     LessonAvatarPositionVector `json:"core_body"`
}

type LessonAvatarPositionVector struct {
  X float32 `json:"x"`
  Y float32 `json:"y"`
  Z float32 `json:"z"`
}
