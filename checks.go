// This file contains all check related functions
package main

import "net"

// This function checks to see if any fields are missing
func checkMissingFields(c *Config) {
	switch {

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
