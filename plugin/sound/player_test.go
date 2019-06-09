package main

import (
	"path"
	"testing"
)

const testVoiceName = "testing_voice"
var audioPath = path.Clean(path.Join(wd, "../../", resourcePath))

func TestExecRandomAlarm(t *testing.T) {
	execRandomAlarm(path.Join(audioPath, alarmPath))
}

func TestExecVoice(t *testing.T) {
	execVoice(path.Join(audioPath, voicePath), testVoiceName)
}
