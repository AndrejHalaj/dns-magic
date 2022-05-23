package main

import (
	"encoding/base64"
	"errors"
	"flag"
	"fmt"

	"github.com/miekg/dns"
)

type EncodeCommand struct {
	fs *flag.FlagSet

	dnsType string
	verbose bool
}

func NewEncodeCommand() *EncodeCommand {
	cmd := &EncodeCommand{
		fs: flag.NewFlagSet("encode", flag.ContinueOnError),
	}

	cmd.fs.StringVar(&cmd.dnsType, "t", "A", "Type of DNS request. (refer to types.go for supported types)")
	cmd.fs.BoolVar(&cmd.verbose, "v", false, "Verbose mode")

	return cmd
}

func (cmd *EncodeCommand) Name() string {
	return cmd.fs.Name()
}

func (cmd *EncodeCommand) Init(args []string) error {
	return cmd.fs.Parse(args)
}

func (cmd *EncodeCommand) Run() error {
	args := cmd.fs.Args()
	if len(args) < 1 {
		return errors.New("missing request hostname")
	}

	hostname := args[0]
	return cmd.encodeAndPrint(cmd.dnsType, hostname)
}

func (cmd *EncodeCommand) encodeAndPrint(t string, hostname string) error {
	tt, err := mapType(t)
	if err != nil {
		return err
	}

	msg := dns.Msg{}
	msg.SetQuestion(normalizeHostname(hostname), tt)

	wire, err := msg.Pack()
	if err != nil {
		return err
	}

	print(base64.StdEncoding.EncodeToString(wire))

	if cmd.verbose {
		println("\nMessage that was encoded:")
		println(msg.String())
	}

	return nil
}

func normalizeHostname(hostname string) string {
	if hostname[len(hostname)-1:] == "." {
		return hostname
	}

	return hostname + "."
}

func mapType(typeStr string) (uint16, error) {
	if t, ok := types[typeStr]; ok {
		return t, nil
	}

	return 0, fmt.Errorf("invalid request type %s, refer to types.go for supported types", typeStr)
}
