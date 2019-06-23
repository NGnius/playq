package controller

type CommandLineInterface struct {
  Channel chan string
}

func NewCLI() CommandLineInterface{
  return CommandLineInterface{}
}
