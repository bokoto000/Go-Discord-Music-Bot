package framework

import (
	"os/exec"
	"strconv"
)

type Song struct {
	Media    string
	Title    string
	Duration *string
	Id       string
	Cheers   string
}

func (song Song) Ffmpeg() *exec.Cmd {
	return exec.Command("ffmpeg", "-i", song.Media, "-f", "s16le", "-ar", strconv.Itoa(FRAME_RATE), "-ac",
		strconv.Itoa(CHANNELS), "pipe:1")
}

func NewSong(media, title, id string) *Song {
	song := new(Song)
	song.Media = media
	song.Title = title
	song.Id = id
	return song
}

func NewSongCheers(media, title, id string, cheers string) *Song {
	song := new(Song)
	song.Media = media
	song.Title = title
	song.Id = id
	song.Cheers = cheers
	return song
}
