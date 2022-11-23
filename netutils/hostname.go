// Package netutils provides useful utilties network specific operations
package netutils

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/Showmax/go-fqdn"
	"github.com/acceldata-io/goutils/shellutils/cmd"
)

const hostnameBinPath = "/bin/hostname"

// GetHostName fetches the hostname of the machine using various methods
// Supported methods are "RFQDN", "FQDN", "OS" and "CMD"
func GetHostName(hostNameCommand string, cmdTimeout int) (string, error) {
	switch hostNameCommand {
	case "RFQDN":
		hostName, err := getReverseFQDN()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(hostName), nil
	case "FQDN":
		hostName, err := fqdn.FqdnHostname()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(hostName), nil
	case "OS":
		hostName, err := os.Hostname()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(hostName), nil
	case "CMD":
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cmdTimeout)*time.Second)
		defer cancel()

		cmdToRun := cmd.New(ctx, "", []string{})
		cmdToRun.WithExpression("bash", hostnameBinPath+" -f")
		cmdToRun.Run()
		if cmdToRun.Status.ExitCode == 0 {
			if strings.TrimSpace(cmdToRun.Status.StdOut) == "" {
				return "", fmt.Errorf("getting empty response in CMD Hostname Method, Because: %s", cmdToRun.Status.StdErr)
			}
			return strings.TrimSpace(cmdToRun.Status.StdOut), nil
		}
		return "", fmt.Errorf(cmdToRun.Status.StdErr)
	default:
		return "", fmt.Errorf("unsupported method %q specified", hostNameCommand)
	}
}

func getReverseFQDN() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	address, err := net.LookupIP(hostname)
	if err != nil {
		return hostname, err
	}

	for _, addr := range address {
		if ipv4 := addr.To4(); ipv4 != nil {
			ip, err := ipv4.MarshalText()
			if err != nil {
				return hostname, err
			}
			hosts, err := net.LookupAddr(string(ip))
			if err != nil || len(hosts) == 0 {
				return hostname, err
			}
			fqdnHostname := hosts[0]
			return strings.TrimSuffix(fqdnHostname, "."), nil // return fqdn without trailing dot
		}
	}
	return hostname, nil
}
