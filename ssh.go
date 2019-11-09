package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSH is the ssh process
func SSH(ip string) {
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
		fmt.Println("*** Authentication Failed - check your credentials. ***")
		fmt.Println()
		return
	}

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	session.Stdout = os.Stdout

	pipe, err := session.StdinPipe()
	if err != nil {
		fmt.Println("Unable to create session pipe: ", err)
		return
	}

	if err = session.Shell(); err != nil {
		fmt.Println("Unable to start shell on server: ", err)
		return
	}

	commands := []string{
		"configure",
	}
	// Run commands
	for _, cmd := range commands {
		fmt.Fprintf(pipe, "%v\n", cmd)
		time.Sleep(1300 * time.Millisecond)
		fmt.Println()
	}

	// Exit config mode
	fmt.Fprint(pipe, "exit\n")
	// Exit session
	fmt.Fprint(pipe, "exit\n")
	// Wait for last command to exit
	if err = session.Wait(); err != nil {
		fmt.Println("The last command did not exit properly: ", err)
		return
	}
}
