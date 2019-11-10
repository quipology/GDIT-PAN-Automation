package main

var pan1Commands []string
var pan2Commands []string

func setPAN1Commands(p *PAN1) []string {
	commands := []string{
		"configure",
		"set deviceconfig system hostname " + p.Hostname,
	}
	return commands
}

func setPAN2Commands(p *PAN2) []string {
	commands := []string{
		"configure",
		"set deviceconfig system hostname " + p.Hostname,
	}
	return commands
}
