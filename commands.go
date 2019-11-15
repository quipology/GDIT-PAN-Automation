// Author: Bobby Williams <quipology@gmail.com>

package main

var pan1Commands []string
var pan2Commands []string

// Commands for PAN1
func setPAN1Commands(p *PAN1) []string {
	commands := []string{
		"configure",
		"set deviceconfig system hostname " + p.Hostname,
	}
	return commands
}

// Commands for PAN2
func setPAN2Commands(p *PAN2) []string {
	commands := []string{
		"configure",
		"set deviceconfig system hostname " + p.Hostname,
	}
	return commands
}
