package main

import (
	"encoding/base64"
	"flag"
	"fmt"

	"github.com/miekg/dns"
)

type EncodeCommand struct {
	fs *flag.FlagSet

	dnsType  string
	hostname string
}

func NewEncodeCommand() *EncodeCommand {
	cmd := &EncodeCommand{
		fs: flag.NewFlagSet("encode", flag.ContinueOnError),
	}

	cmd.fs.StringVar(&cmd.dnsType, "t", "A", "Type of DNS request. (A/AAAA)")
	cmd.fs.StringVar(&cmd.hostname, "host", "example.com", "Hostname of the DNS request.")

	return cmd
}

func (cmd *EncodeCommand) Name() string {
	return cmd.fs.Name()
}

func (cmd *EncodeCommand) Init(args []string) error {
	return cmd.fs.Parse(args)
}

func (cmd *EncodeCommand) Run() error {
	encoded, err := encode(cmd.dnsType, cmd.hostname)
	if err != nil {
		return err
	} else {
		println(encoded)
	}

	return nil
}

func encode(t string, hostname string) (string, error) {
	tt, err := mapType(t)
	if err != nil {
		return "", err
	}

	q := dns.Question{
		Name:   normalizeHostname(hostname),
		Qtype:  tt,
		Qclass: dns.ClassINET,
	}

	msg := dns.Msg{}
	msg.Question = append(msg.Question, q)

	wire, err := msg.Pack()
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(wire), nil
}

func normalizeHostname(hostname string) string {
	if hostname[len(hostname)-1:] == "." {
		return hostname
	}

	return hostname + "."
}

func mapType(t string) (uint16, error) {
	switch t {
	case "A":
		return dns.TypeA, nil
	case "AAAA":
		return dns.TypeAAAA, nil
	case "SVCB":
		return dns.TypeSVCB, nil
	default:
		return 0, fmt.Errorf("Invalid request type %s. Only A/AAAA/SVCB are supported.", t)
	}
}
