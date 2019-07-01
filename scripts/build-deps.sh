#!/bin/bash
echo "NOTE: You may need to install libasound2-dev or openal first"
echo "(more info: https://github.com/hajimehoshi/oto/blob/master/README.md)"
# install dependencies
go get -u github.com/faiface/beep github.com/gdamore/tcell github.com/hajimehoshi/go-mp3 github.com/hajimehoshi/oto	github.com/jfreymuth/oggvorbis github.com/jfreymuth/vorbis github.com/mewkiz/flac github.com/pkg/errors
if [ $? != 0 ]
then
  echo "ERROR: Dependencies failed to install"
  exit $?
fi
echo "Dependencies successfully installed"
# actual build
go build src/main/playq.go
if [ $? != 0 ]
then
  echo "ERROR: Build failed"
  exit $?
fi
echo "Build successful, all dependencies satisfied"
echo "(Use br.sh for faster build & run from now on)"
