package model

type PlayerStatus struct {
	IsPlaying bool `json:"is_playing"`
	Volume    int  `json:"volume"`
}
