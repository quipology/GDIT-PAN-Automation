package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	MgntIP             string
	Hostname           string
	TunnelName         string
	TunnelInterface    string
	VirtualRouter      string
	IKEprofile         string
	IKEgateway         string
	IPSECprofile       string
	PeerIP             string
	InterestingTraffic []string
}

const (
	username = "admin"
	password = "admin"
	port     = ":22"

	logFile = "pancfg.log"
)

// Logging function to set destination log file
func setLogger() *os.File {
	l, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Panic(err)
	}
	return l
}

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

func main() {
	// -----------------
	// | Setup logging |
	// -----------------
	// Generate unique session ID for logging
	rand.Seed(time.Now().UnixNano())
	sessionID := rand.Intn(99999)
	// Set log file
	log.SetOutput(setLogger())
	// Start logging
	log.Printf("**** Session - %v - started ****\n", sessionID)
	// --------------------
	// | Read config file |
	// -------------------
	// Config filename flag (default: config.yml)
	configFile := flag.String("f", "config.yml", "Specify your config `filename` (ex. pancfg -f config.yml)")
	flag.Parse()
	// Check if config file is *.yml
	if !strings.HasSuffix(strings.ToLower(*configFile), ".yml") {
		fmt.Println("error: config must be YAML file (.yml)")
	}
	// Open config file
	file, err := os.Open(*configFile)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	// Read file
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// Decode YAML
	var c Config
	if err = yaml.Unmarshal(fileBytes, &c); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Loading config file '%v'..\n", *configFile)
	// Check for missing fields in config file
	fmt.Println("Checking for missing fields..")
	checkMissingFields(&c)
}
