package playback

import (
    //"fmt"
    //"flag"
    //"github.com/faiface/beep"
    "os"
)

type Playback struct {
  filename string
  currentFile os.File
  Channel chan os.File
  Events chan string
}

func (p Playback) Start(end chan int) {
  go p.Run(end)
}

func (p Playback) Run(end chan int) {
  end <- 0
}

func New(channel chan os.File) Playback{
  return Playback{Channel:channel, Events:make(chan string)}
}
