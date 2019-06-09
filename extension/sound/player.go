package main

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	"github.com/LinMAD/BitAccretion/extension"
	"github.com/LinMAD/BitAccretion/model"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

const (
	resourcePath = "/resource/sound"
	alarmPath    = "alarm"
	voicePath    = "voice"
)

// Player for sound making
type Player struct {
	// wd is work dir or rooted base path
	wd string
}

// PlayAlert for given name
func (p *Player) PlayAlert(name model.VertexName) {
	if p.wd == "" {
		return
	}

	execRandomAlarm(path.Join(p.wd, resourcePath, alarmPath))
	execVoice(path.Join(p.wd, resourcePath, voicePath), string(name))
}

// NewSound prepares sound player
func NewSound() extension.ISound {
	rand.Seed(time.Now().UTC().UnixNano())
	wd, _ := os.Getwd()

	return &Player{wd: wd}
}

// play given streamer and rate of it
func play(s beep.Streamer, r beep.SampleRate) {
	_ = speaker.Init(r, r.N(time.Second/10))

	playing := make(chan struct{})
	speaker.Play(beep.Seq(s, beep.Callback(func() {
		close(playing)
	})))

	<-playing
}

// execRandomAlarm sound file
func execRandomAlarm(dir string) {
	files, filesErr := ioutil.ReadDir(dir)
	if filesErr != nil {
		return
	}

	name := files[rand.Intn(len(files))].Name()
	af, _ := os.Open(path.Join(dir, name))
	stream, format, _ := wav.Decode(af)

	play(stream, format.SampleRate)
}

// execVoice matching voice name if exist
func execVoice(dir, reqName string) {
	files, filesErr := ioutil.ReadDir(dir)
	if filesErr != nil {
		return
	}

	var voiceFileName string
	for _, f := range files {
		if strings.Contains(f.Name(), reqName) {
			voiceFileName = f.Name()
			break
		}
	}

	if voiceFileName == "" {
		return
	}

	audioFile, _ := os.Open(path.Join(dir, voiceFileName))
	stream, format, _ := wav.Decode(audioFile)
	play(stream, format.SampleRate)
}
