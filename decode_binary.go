package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/miekg/dns"
)

type DecodeBinaryCommand struct {
	fs *flag.FlagSet
}

func NewDecodeBinaryCommand() *DecodeBinaryCommand {
	cmd := &DecodeBinaryCommand{
		fs: flag.NewFlagSet("decode-binary", flag.ContinueOnError),
	}

	return cmd
}

func (cmd *DecodeBinaryCommand) Name() string {
	return cmd.fs.Name()
}

func (cmd *DecodeBinaryCommand) Init(args []string) error {
	return cmd.fs.Parse(args)
}

func (cmd *DecodeBinaryCommand) Run() error {
	msgData, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("error reading input: %v", err)
	}

	var msg dns.Msg
	err = msg.Unpack(msgData)
	if err != nil {
		return fmt.Errorf("error unpacking DNS message: %v", err)
	}

	print(msg.String())
	return nil
}
