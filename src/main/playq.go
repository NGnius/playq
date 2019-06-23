package main

import (
    "fmt"
    // "flag"
    "os"
    "../controller"
    "../monitor"
    "../playback"
)

func main() {
  // create CLI
  cli := controller.NewCLI()

  // create monitor
  mon := monitor.New("A")
  cli.Channel = mon.ControlChannel

  // create player
  playback_complete := make(chan int)
  player := playback.New(make(chan os.File))

  player.Start(playback_complete)

  fmt.Println(<- playback_complete)
}
