package main

type mode int

const (
	normal mode = iota
	insert
	command
)

func (m mode) String() string {
	switch m {
	case insert:
		return "INS"
	case command:
		return "CMD"
	default:
		return "NRM"
	}
}
