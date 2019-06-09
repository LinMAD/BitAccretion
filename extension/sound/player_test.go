package main

import (
	"os"
	"path"
	"testing"
)

const testVoiceName = "testing_voice"

var (
	wd        string
	audioPath = path.Clean(path.Join(wd, "../../", resourcePath))
)

func init() {
	wd, _ = os.Getwd()
}

func TestExecRandomAlarm(t *testing.T) {
	execRandomAlarm(path.Join(audioPath, alarmPath))
}

func TestExecVoice(t *testing.T) {
	execVoice(path.Join(audioPath, voicePath), testVoiceName)
}
