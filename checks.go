// Author: Bobby Williams <quipology@gmail.com>

package main

import (
	"fmt"
	"net"
)

// This function checks to see if any fields are missing
func checkMissingFieldsPAN1(p *PAN1) {
	switch {
	case p.MgntIP == "":
		fmt.Print(errMissingConfigField, ": pan1-mgnt-ip\n")
		displayExit()
	case p.Hostname == "":
		fmt.Print(errMissingConfigField, ": pan1-hostname\n")
		displayExit()
	case p.TunnelName == "":
		fmt.Print(errMissingConfigField, ": pan1-tunnel-name\n")
		displayExit()
	case p.TunnelIPMask == "":
		fmt.Print(errMissingConfigField, ": pan1-tunnel-ip-and-mask\n")
		displayExit()
	case p.VirtualRouter == "":
		fmt.Print(errMissingConfigField, ": pan1-virtual-router\n")
		displayExit()
	case p.IKEprofile == "":
		fmt.Print(errMissingConfigField, ": pan1-ike-crypto-profile\n")
		displayExit()
	case p.IKEgateway == "":
		fmt.Print(errMissingConfigField, ": pan1-ike-gateway\n")
		displayExit()
	case p.IPSECprofile == "":
		fmt.Print(errMissingConfigField, ": pan1-ipsec-crypto-profile\n")
		displayExit()
	case p.PeerIPMask == "":
		fmt.Print(errMissingConfigField, ": pan1-peer-ip-and-mask\n")
		displayExit()
	case len(p.InterestingTraffic) == 0:
		fmt.Print(errMissingConfigField, ": pan1-destination-subnets\n")
		displayExit()
	}
}

// This function checks to see if any fields are missing
func checkMissingFieldsPAN2(p *PAN2) {
	switch {
	case p.MgntIP == "":
		fmt.Print(errMissingConfigField, ": pan2-mgnt-ip\n")
		displayExit()
	case p.Hostname == "":
		fmt.Print(errMissingConfigField, ": pan2-hostname\n")
		displayExit()
	case p.TunnelName == "":
		fmt.Print(errMissingConfigField, ": pan2-tunnel-name\n")
		displayExit()
	case p.TunnelIPMask == "":
		fmt.Print(errMissingConfigField, ": pan2-tunnel-ip-and-mask\n")
		displayExit()
	case p.VirtualRouter == "":
		fmt.Print(errMissingConfigField, ": pan2-virtual-router\n")
		displayExit()
	case p.IKEprofile == "":
		fmt.Print(errMissingConfigField, ": pan2-ike-crypto-profile\n")
		displayExit()
	case p.IKEgateway == "":
		fmt.Print(errMissingConfigField, ": pan2-ike-gateway\n")
		displayExit()
	case p.IPSECprofile == "":
		fmt.Print(errMissingConfigField, ": pan2-ipsec-crypto-profile\n")
		displayExit()
	case p.PeerIPMask == "":
		fmt.Print(errMissingConfigField, ": pan2-peer-ip-and-mask\n")
		displayExit()
	case len(p.InterestingTraffic) == 0:
		fmt.Print(errMissingConfigField, ": pan2-destination-subnets\n")
		displayExit()
	}
}

// This function checks to see if passed in IP is valid IPv4 format
func checkIP(s string) bool {
	if check := net.ParseIP(s); check == nil {
		return false
	}
	return true
}

// This function checks to see if passed in IPMask is valid CIDR notation
func checkIPMask(s string) error {
	if _, _, err := net.ParseCIDR(s); err != nil {
		return err
	}
	return nil
}

// This function checks for a valid tunnel number
func checkTunnelNumber(n uint32, name string) {
	if n == 0 {
		displayOutput(fmt.Sprintf("\t - Error: %v-tunnel-number can't be 0.", name))
		displayExit()
	}
}
