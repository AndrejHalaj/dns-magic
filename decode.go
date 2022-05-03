package main

import (
	"encoding/base64"
	"flag"

	"github.com/miekg/dns"
)

type DecodeCommand struct {
	fs *flag.FlagSet

	dnsRequest string
}

func NewDecodeCommand() *DecodeCommand {
	cmd := &DecodeCommand{
		fs: flag.NewFlagSet("decode", flag.ContinueOnError),
	}

	cmd.fs.StringVar(&cmd.dnsRequest, "request", "", "DNS request in wireformat.")

	return cmd
}

func (cmd *DecodeCommand) Name() string {
	return cmd.fs.Name()
}

func (cmd *DecodeCommand) Init(args []string) error {
	return cmd.fs.Parse(args)
}

func (cmd *DecodeCommand) Run() error {
	decoded, err := decode(cmd.dnsRequest)
	if err != nil {
		return err
	} else {
		println(decoded)
	}

	return nil
}

func decode(dnsMsg string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(dnsMsg)
	if err != nil {
		return "", err
	}

	msg := dns.Msg{}
	err = msg.Unpack(decoded)
	if err != nil {
		return "", err
	}

	return msg.String(), nil
}
