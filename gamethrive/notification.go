package gamethrive

import (
	"time"
)

type NotificationsService struct {
	c *Client
}

type Notification struct {
	Id string `json:"-"`
	// Required Body Parameters
	AppId     string            `json:"app_id"`
	IsIOS     bool              `json:"isIos"`
	IsAndroid bool              `json:"isAndroid"`
	Contents  map[string]string `json:"contents"`
	// Target Parameters
	IncludedSegments      []string `json:"included_segments,omitempty"`
	ExcludedSegments      []string `json:"excluded_segments,omitempty"`
	IncludedPlayerIds     []string `json:"include_player_ids,omitempty"`
	IncludedIOSTokens     []string `json:"include_ios_tokens,omitempty"`
	IncludedAndroidRegIds []string `json:"include_android_reg_ids,omitempty"`
	// Optional Body Paramters

	ContentAvailable   bool              `json:"content_available,omitempty"`
	IOSBadgeType       BadgeType         `json:"ios_badgeType,omitempty"`
	IOSBadgeCount      int               `json:"ios_badgeCount,omitempty"`
	IOSSound           string            `json:"ios_sound,omitempty"`
	AndroidSound       string            `json:"android_sound,omitempty"`
	Data               map[string]string `json:"data,omitempty"`
	URL                string            `json:"url,omitempty"`
	SendAfter          *time.Time        `json:"send_after,omitempty"`
	SendUserActiveTime bool              `json:"send_at_user_active_time,omitempty"`
}

type BadgeType string

const (
	None     BadgeType = "None"
	SetTo    BadgeType = "SetTo"
	Increase BadgeType = "Increase"
)

func (s *NotificationsService) New(notification *Notification, auth string) (int, error) {
	req, err := s.c.NewRequest("POST", "notifications", notification)
	if err != nil {
		return 0, err
	}
	if len(auth) > 0 {
		req.Header.Set("Authorization", "Basic "+auth)
	}
	var res struct {
		Id         string `json:"id"`
		Recipients int    `json:"recipients"`
	}
	_, err = s.c.Do(req, &res)
	if err != nil {
		return 0, err
	}
	notification.Id = res.Id
	return res.Recipients, nil
}

func (s *NotificationsService) Open(notification *Notification, opened bool) error {
	urlStr := "notifications/" + notification.Id
	body := struct {
		Opened bool   `json:"opened"`
		AppId  string `json:"app_id"`
	}{
		Opened: opened,
		AppId:  notification.AppId,
	}
	req, err := s.c.NewRequest("PUT", urlStr, body)
	if err != nil {
		return err
	}
	_, err = s.c.Do(req, nil)
	return err
}
