package gamethrive

import (
	"errors"
	"fmt"
)

type PlayersService struct {
	c *Client
}

type Player struct {
	// Required fields
	DeviceType DeviceType `json:"device_type"`
	AppId      string     `json:"app_id"`
	Id         string     `json:"id,omitempty"`
	// Recommended fields
	Identifier    string `json:"identifier,omitempty"`
	Language      string `json:"language,omitempty"`
	Timezone      int    `json:"timezone,omitempty"`
	DeviceModel   string `json:"device_model,omitempty"`
	DeviceOS      string `json:"device_os,omitempty"`
	GameVersion   string `json:"game_version,omitempty"`
	AdvertisingId string `json:"ad_id,omitempty"`
	// Other/Optional fields
	SessionCount int               `json:"session_count,omitempty"`
	Tags         map[string]string `json:"tags,omitempty"`
	AmountSpent  float64           `json:"amount_spent,omitempty"`
	CreatedAt    int               `json:"created_at,omitempty"`
	LastActive   int               `json:"last_active,omitempty"`
	Playtime     int               `json:"playtime,omitempty"`
}

type DeviceType int

const (
	IOS     DeviceType = 0
	Android DeviceType = 1
	Amazon  DeviceType = 2
)

type PlaytimeState string

const (
	Suspend PlaytimeState = "suspend"
	Resume  PlaytimeState = "resume"
	Ping    PlaytimeState = "ping"
)

// Todo: test player id
func (s *PlayersService) New(player *Player) error {
	req, err := s.c.NewRequest("POST", "/players", player)
	if err != nil {
		return err
	}
	var res struct {
		Success bool   `json:"success"`
		Id      string `json:"id"`
	}
	_, err = s.c.Do(req, &res)
	if err != nil {
		return err
	}
	player.Id = res.Id
	return nil
}

func (s *PlayersService) Update(player *Player) error {
	if len(player.Id) <= 0 {
		return errors.New("Player id is required")
	}
	urlStr := fmt.Sprintf("players/%s", player.Id)
	req, err := s.c.NewRequest("PUT", urlStr, player)
	if err != nil {
		return err
	}
	_, err = s.c.Do(req, nil)
	return err
}

func (s *PlayersService) UpdateAmount(playerId string, amount float64) error {
	if len(playerId) <= 0 {
		return errors.New("Player id is required")
	}
	urlStr := fmt.Sprintf("players/%s/on_purchase", playerId)
	body := struct {
		Amount float64 `json:"amount"`
	}{
		Amount: amount,
	}
	req, err := s.c.NewRequest("POST", urlStr, body)
	if err != nil {
		return err
	}
	_, err = s.c.Do(req, nil)
	return err
}

func (s *PlayersService) Session(player *Player) error {
	if len(player.Id) <= 0 {
		return errors.New("Player id is required")
	}
	urlStr := fmt.Sprintf("players/%s/on_session", player.Id)
	req, err := s.c.NewRequest("POST", urlStr, player)
	if err != nil {
		return err
	}
	_, err = s.c.Do(req, nil)
	return err
}

func (s *PlayersService) Playtime(playerId string, state PlaytimeState, time int) error {
	if len(playerId) <= 0 {
		return errors.New("Player id is required")
	}
	urlStr := fmt.Sprintf("players/%s/on_focus", playerId)
	body := struct {
		State string `json:"state"`
		Time  int    `json:"active_time"`
	}{
		State: string(state),
		Time:  time,
	}
	req, err := s.c.NewRequest("POST", urlStr, body)
	if err != nil {
		return err
	}
	_, err = s.c.Do(req, nil)
	return err
}
