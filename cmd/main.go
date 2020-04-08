package main

import (
	"fmt"
	"github.com/code-ready/goodhosts"
	"github.com/docopt/docopt-go"
	"net"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	usage := `Goodhosts - simple hosts file management.

Usage:
  goodhosts check <ip> <host>...
  goodhosts add <ip> <host>...
  goodhosts list [--all]
  goodhosts (rm|remove) <ip> | <host>...
  goodhosts -h | --help
  goodhosts --version

Options:
  --all         Display comments when listing.
  -h --help     Show this screen.
  --version     Show the version.`

	args, _ := docopt.Parse(usage, nil, true, "Goodhosts 2.2.2", false)

	hosts, err := goodhosts.NewHosts()
	check(err)

	if args["list"].(bool) {
		total := 0
		for _, line := range hosts.Lines {
			var lineOutput string

			if line.IsComment() && !args["--all"].(bool) {
				continue
			}

			lineOutput = fmt.Sprintf("%s", line.Raw)
			if line.Err != nil {
				lineOutput = fmt.Sprintf("%s # <<< Malformated!", lineOutput)
			}
			total += 1

			fmt.Println(lineOutput)
		}

		fmt.Printf("\nTotal: %d\n", total)

		return
	}

	if args["check"].(bool) {
		hasErr := false

		ip := args["<ip>"].(string)
		hostEntries := args["<host>"].([]string)

		for _, hostEntry := range hostEntries {
			if !hosts.Has(ip, hostEntry) {
				fmt.Fprintln(os.Stderr, fmt.Sprintf("%s %s is not in the hosts file", ip, hostEntry))
				hasErr = true
			}
		}

		if hasErr {
			os.Exit(1)
		}

		return
	}

	if args["add"].(bool) {
		ip := args["<ip>"].(string)
		hostEntries := args["<host>"].([]string)

		if !hosts.IsWritable() {
			fmt.Fprintln(os.Stderr, "Host file not writable. Try running with elevated privileges.")
			os.Exit(1)
		}

		err = hosts.Add(ip, hostEntries...)
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("%s", err.Error()))
			os.Exit(2)
		}

		err = hosts.Flush()
		check(err)

		return
	}

	if args["rm"].(bool) || args["remove"].(bool) {
		ip := args["<ip>"].(string)
		hostEntries := args["<host>"].([]string)

		arg := []string{}
		arg = append(arg, ip)
		arg = append(arg, hostEntries...)

		if !hosts.IsWritable() {
			fmt.Fprintln(os.Stderr, "Host file not writable. Try running with elevated privileges.")
			os.Exit(1)
		}

		err = remove(&hosts, arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", err.Error()))
			os.Exit(2)
		}

		err = hosts.Flush()
		check(err)

		return
	}
}

func remove(hosts *goodhosts.Hosts, arg []string) error {
	if len(arg) == 0 {
		return fmt.Errorf("Not enough arguments")
	}

	if len(arg) == 1 {
		fmt.Println("Processing single arg")
		processSingleArg(hosts, arg[0])
	}

	uniqueHosts := map[string]bool{}
	var hostEntries []string

	for i := 1; i < len(arg); i++ {
		uniqueHosts[arg[i]] = true
	}

	for key, _ := range uniqueHosts {
		hostEntries = append(hostEntries, key)
	}

	if net.ParseIP(arg[0]) != nil {
		if hosts.HasIp(arg[0]) {
			err := hosts.Remove(arg[0], hostEntries...)
			if err != nil {
				return err
			}
		}
	} else {
		hostEntries = append(hostEntries, arg[0])
		for _, value := range hostEntries {
			err := hosts.RemoveByHostname(value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func processSingleArg(host *goodhosts.Hosts, arg string) error {
	if net.ParseIP(arg) != nil {
		fmt.Println("Removing using IP")
		if err := host.RemoveByIp(arg); err != nil {
			return err
		}
		return nil
	}

	if err := host.RemoveByHostname(arg); err != nil {
		return err
	}
	return nil
}