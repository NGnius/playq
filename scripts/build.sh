#!/bin/bash
# install dependencies
go get -u github.com/faiface/beep github.com/gdamore/tcell github.com/hajimehoshi/go-mp3 github.com/hajimehoshi/oto	github.com/jfreymuth/oggvorbis github.com/jfreymuth/vorbis github.com/mewkiz/flac github.com/pkg/errors
# actual build
go build src/main/playq.go
