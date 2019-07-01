package main

import (
    "fmt"
    // "flag"
    "../controller"
    "../monitor"
    "../playback"
)

var ABOUT_STRING string = `PlayQ v0.0.1 for StreamQ
Developed by NGnius (Graham Littlewood)
---------------------------------------`

func main() {
  fmt.Println(ABOUT_STRING)
  // create CLI
  cli := controller.NewCLI()
  cli_complete := make(chan int)

  // create monitor
  mon := monitor.New("A")
  cli.MonitorControlChannel = mon.ControlChannel
  cli.MonitorEventChannel = mon.EventChannel
  monitor_complete := make(chan int)

  // create player
  player := playback.New()
  mon.PlaybackControlChannel = player.ControlChannel
  mon.PlaybackFileChannel = player.FileChannel
  mon.PlaybackEventChannel = player.EventChannel
  playback_complete := make(chan int)

  // start components
  cli.Start(cli_complete)
  player.Start(playback_complete)
  mon.Start(monitor_complete)

  // wait until all components complete
  //fmt.Println(<- cli_complete)
  <- cli_complete
  //fmt.Println(<- monitor_complete)
  <- monitor_complete
  //fmt.Println(<- playback_complete)
  <- playback_complete
}
