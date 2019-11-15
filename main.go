// Author: Bobby Williams <quipology@gmail.com>

package main

import (
	"bufio"
	crand "crypto/rand"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	yaml "gopkg.in/yaml.v2"
)

const (
	username = "admin" // SSH username
	password = "admin" // SSH password
	port     = ":22"   // SSH port

	preSharedKeyLength = 30 // Desired length of pre-shared-key

	logFile = "pancfg.log" // Log file to be created/appended
)

var (
	preSharedKey string // Variable to store generated pre-shared-key

	errMissingConfigField = errors.New("Missing one or more required config fields")
	errInvalidIPv4        = errors.New("Invalid IPv4 address")
	errInvalidIPMask      = errors.New("Invalid CIDR")

	wg sync.WaitGroup
)

// PAN1 represents Firewall 1's config settings
type PAN1 struct {
	MgntIP             string   `yaml:"pan1-mgnt-ip"`
	Hostname           string   `yaml:"pan1-hostname"`
	TunnelName         string   `yaml:"pan1-tunnel-name"`
	TunnelNumber       uint32   `yaml:"pan1-tunnel-number"`
	TunnelIPMask       string   `yaml:"pan1-tunnel-ip-and-mask"`
	LocalIPMask        string   `yaml:"pan1-local-ip-and-mask"`
	VirtualRouter      string   `yaml:"pan1-virtual-router"`
	IKEprofile         string   `yaml:"pan1-ike-crypto-profile"`
	IKEgateway         string   `yaml:"pan1-ike-gateway"`
	IPSECprofile       string   `yaml:"pan1-ipsec-crypto-profile"`
	PeerIPMask         string   `yaml:"pan1-peer-ip-and-mask"`
	InterestingTraffic []string `yaml:"pan1-destination-subnets"`
	PSK                string   // Pre-shared-key
}

// PAN2 represents Firewall 2's config settings
type PAN2 struct {
	MgntIP             string   `yaml:"pan2-mgnt-ip"`
	Hostname           string   `yaml:"pan2-hostname"`
	TunnelName         string   `yaml:"pan2-tunnel-name"`
	TunnelNumber       uint32   `yaml:"pan2-tunnel-number"`
	TunnelIPMask       string   `yaml:"pan2-tunnel-ip-and-mask"`
	LocalIPMask        string   `yaml:"pan2-local-ip-and-mask"`
	VirtualRouter      string   `yaml:"pan2-virtual-router"`
	IKEprofile         string   `yaml:"pan2-ike-crypto-profile"`
	IKEgateway         string   `yaml:"pan2-ike-gateway"`
	IPSECprofile       string   `yaml:"pan2-ipsec-crypto-profile"`
	PeerIPMask         string   `yaml:"pan2-peer-ip-and-mask"`
	InterestingTraffic []string `yaml:"pan2-destination-subnets"`
	PSK                string   // Pre-shared-key
}

// Logging function to set destination log file
func setLogger() *os.File {
	l, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		displayOutput(err.Error())
		displayExit()
	}
	return l
}

// Display exit for errors
func displayExit() {
	fmt.Println("Press enter to exit.")
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		os.Exit(1)
	}
}

// For printing to terminal + log
func displayOutput(s string) {
	log.Println(s)
	fmt.Println(s)
}

func main() {
	// Setup Logging and Generate unique session ID
	rand.Seed(time.Now().UnixNano())
	sessionID := rand.Intn(99999)
	// Set log file
	log.SetOutput(setLogger())
	// Start logging
	log.Printf("**** Session - %v - start ****\n", sessionID)

	//--------------------------------------------------
	// Config filename flag (default: config.yml)
	configFile := flag.String("f", "config.yml", "Specify your config `filename` (ex. pancfg -f config.yml)")
	flag.Parse()
	// Check if config file is *.yml
	if !strings.HasSuffix(strings.ToLower(*configFile), ".yml") && !strings.HasSuffix(strings.ToLower(*configFile), ".yaml") {
		displayOutput("error: config must be YAML file (.yml or .yaml)")
	}
	// Read config file
	fileBytes, err := ioutil.ReadFile(*configFile)
	if err != nil {
		displayOutput(err.Error())
		displayExit()
	}
	fmt.Printf("Loading config file '%v'..\n", *configFile)
	//--------------------------------------------------
	// Decode YAML file
	var pan1 PAN1
	if err = yaml.Unmarshal(fileBytes, &pan1); err != nil {
		displayOutput(err.Error())
		os.Exit(1)
	}
	var pan2 PAN2
	if err = yaml.Unmarshal(fileBytes, &pan2); err != nil {
		displayOutput(err.Error())
		os.Exit(1)
	}

	// Generate pre-shared-key for IPSEC tunnel
	shell := make([]byte, preSharedKeyLength)
	if _, err = crand.Read(shell); err != nil {
		displayOutput(err.Error())
		displayExit()
	}
	preSharedKey = base64.URLEncoding.EncodeToString(shell)
	pan1.PSK = preSharedKey
	pan2.PSK = preSharedKey

	// Check for missing fields in config file
	displayOutput("Checking for missing fields..")
	checkMissingFieldsPAN1(&pan1)
	checkMissingFieldsPAN2(&pan2)

	// Check for any invalid IPv4
	displayOutput("Checking for invalid IPv4 syntax..")
	IPfields := []string{
		pan1.MgntIP,
		pan2.MgntIP,
	}
	for _, i := range IPfields {
		if pass := checkIP(i); !pass {
			fmt.Println(errInvalidIPv4)
			displayExit()
		}
	}

	// Check for any invalid IPMask
	displayOutput("Checking for invalid IPMask syntax..")
	IPMaskFields := []string{
		pan1.TunnelIPMask,
		pan2.TunnelIPMask,
		pan1.LocalIPMask,
		pan2.LocalIPMask,
		pan1.PeerIPMask,
		pan2.PeerIPMask,
	}

	wg.Add(3)
	go func() {
		for _, i := range IPMaskFields {
			if err := checkIPMask(i); err != nil {
				displayOutput(err.Error())
				displayExit()
			}
		}
		wg.Done()
	}()

	go func() {
		for _, i := range pan1.InterestingTraffic {
			if err := checkIPMask(i); err != nil {
				displayOutput(err.Error())
				displayExit()
			}
		}
		wg.Done()
	}()

	go func() {
		for _, i := range pan2.InterestingTraffic {
			if err := checkIPMask(i); err != nil {
				displayOutput(err.Error())
				displayExit()
			}
		}
		wg.Done()
	}()

	wg.Wait()
	// Check for valid tunnel number
	displayOutput("Checking for invalid Tunnel number..")
	checkTunnelNumber(pan1.TunnelNumber, "pan1")
	checkTunnelNumber(pan2.TunnelNumber, "pan2")

	// Notify that all checks passed if arrived this far
	displayOutput("All Checks Passed!")

	// Get PAN1 & PAN2 Commands
	pan1Commands = setPAN1Commands(&pan1)
	pan2Commands = setPAN2Commands(&pan2)

	// SSH into PANs to configure
	wg.Add(2)
	// Configure PAN1
	go func() {
		displayOutput("Connecting to PAN1..")
		SSH(pan1.MgntIP, "PAN1", pan1Commands)
		displayOutput("PAN1 configuration complete.")
		wg.Done()
	}()
	go func() {
		displayOutput("Connecting to PAN2..")
		SSH(pan2.MgntIP, "PAN2", pan2Commands)
		displayOutput("PAN2 configuration complete.")
		wg.Done()
	}()

	wg.Wait() // Wait for PANs to be configured

	// End logging
	log.Printf("**** Session - %v - end ****\n", sessionID)
}
