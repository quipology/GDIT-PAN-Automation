package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSH is the ssh process
func SSH(ip, name string, commands []string) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", ip+port, config)
	if err != nil {
		fmt.Println()
		displayOutput("*** Authentication Failed - check your credentials. ***")
		fmt.Println()
		return
	}
	displayOutput(fmt.Sprintf("Login to %v(%v) successful.", name, ip))

	session, err := client.NewSession()
	if err != nil {
		displayOutput("Failed to create session: " + err.Error())
		displayExit()
	}
	defer session.Close()

	pipe, err := session.StdinPipe()
	if err != nil {
		displayOutput("Unable to create session pipe: " + err.Error())
		displayExit()
	}

	var b bytes.Buffer
	session.Stdout = &b
	// session.Stdout = os.Stdout

	if err = session.Shell(); err != nil {
		fmt.Println("Unable to start shell on server: ", err)
		return
	}

	// Sleep a few seconds to allow the PAN CLI to load
	time.Sleep(5 * time.Second)

	// Run commands
	for _, cmd := range commands {
		fmt.Printf("%v: %v\n", name, cmd)
		fmt.Fprintf(pipe, "%v\n", cmd)
		time.Sleep(1300 * time.Millisecond)
	}

	// Commit PAN changes
	b.Reset() // Empty buffer so that commit results can be shown
	fmt.Printf("%v: %v\n", name, "commit")
	fmt.Fprint(pipe, "commit\n")

	// Wait for changes to be committed before exiting
	fmt.Printf("Waiting for %v changes to commit..\n", name)
	for {
		if strings.Contains(b.String(), "100%") {
			break
		}
	}

	// Exit the PAN CLI
	fmt.Fprint(pipe, "exit\n")
	fmt.Fprint(pipe, "exit\n")

}
