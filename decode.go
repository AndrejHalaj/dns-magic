package main

import (
	"encoding/base64"
	"errors"
	"flag"

	"github.com/miekg/dns"
)

type DecodeCommand struct {
	fs *flag.FlagSet
}

func NewDecodeCommand() *DecodeCommand {
	cmd := &DecodeCommand{
		fs: flag.NewFlagSet("decode", flag.ContinueOnError),
	}

	return cmd
}

func (cmd *DecodeCommand) Name() string {
	return cmd.fs.Name()
}

func (cmd *DecodeCommand) Init(args []string) error {
	return cmd.fs.Parse(args)
}

func (cmd *DecodeCommand) Run() error {
	args := cmd.fs.Args()
	if len(args) < 1 {
		return errors.New("missing request to decode")
	}

	dnsRequest := args[0]
	return cmd.decodeAndPrint(dnsRequest)
}

func (cmd *DecodeCommand) decodeAndPrint(dnsMsg string) error {
	decoded, err := base64.StdEncoding.DecodeString(dnsMsg)
	if err != nil {
		return err
	}

	msg := dns.Msg{}
	err = msg.Unpack(decoded)
	if err != nil {
		return err
	}

	print(msg.String())

	return nil
}
