package main

import (
	crand "crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
)

const (
	username = "admin"
	password = "admin"
	port     = ":22"

	preSharedKeyLength = 30 // Desired length of pre-shared-key

	logFile = "pancfg.log"
)

var preSharedKey string

// PAN1 represents Firewall 1's config settings
type PAN1 struct {
	MgntIP             string   `yaml:"pan1-mgnt-ip`
	Hostname           string   `yaml:"pan1-hostname"`
	TunnelName         string   `yaml:"pan1-tunnel-name"`
	TunnelInterface    string   `yaml:"pan1-tunnel-interface`
	TunnelIPMask       string   `yaml:"pan1-tunnel-ip-and-mask"`
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
	MgntIP             string   `yaml:"pan2-mgnt-ip`
	Hostname           string   `yaml:"pan2-hostname"`
	TunnelName         string   `yaml:"pan2-tunnel-name"`
	TunnelInterface    string   `yaml:"pan2-tunnel-interface`
	TunnelIPMask       string   `yaml:"pan2-tunnel-ip-and-mask"`
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
		log.Panic(err)
	}
	return l
}

func main() {
	// Setup Logging and Generate unique session ID
	rand.Seed(time.Now().UnixNano())
	sessionID := rand.Intn(99999)
	// Set log file
	log.SetOutput(setLogger())
	// Start logging
	log.Printf("**** Session - %v - started ****\n", sessionID)

	//--------------------------------------------------
	// Config filename flag (default: config.yml)
	configFile := flag.String("f", "config.yml", "Specify your config `filename` (ex. pancfg -f config.yml)")
	flag.Parse()
	// Check if config file is *.yml
	if !strings.HasSuffix(strings.ToLower(*configFile), ".yml") {
		fmt.Println("error: config must be YAML file (.yml)")
	}
	// Read config file
	fileBytes, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Panic(err)
	}
	//--------------------------------------------------
	// Decode YAML file
	var pan1 PAN1
	if err = yaml.Unmarshal(fileBytes, &pan1); err != nil {
		fmt.Println(err)
	}
	var pan2 PAN2
	if err = yaml.Unmarshal(fileBytes, &pan2); err != nil {
		fmt.Println(err)
	}

	// Generate pre-shared-key for IPSEC tunnel
	shell := make([]byte, preSharedKeyLength)
	if _, err = crand.Read(shell); err != nil {
		log.Panic(err)
	}
	preSharedKey = base64.URLEncoding.EncodeToString(shell)
	pan1.PSK = preSharedKey
	pan2.PSK = preSharedKey

	fmt.Printf("Loading config file '%v'..\n", *configFile)
	// Check for missing fields in config file
	fmt.Println("Checking for missing fields..")
	checkMissingFields(&pan1)

	// Start logging
	log.Printf("**** Session - %v - finished ****\n", sessionID)
}
