package main

import (
	"bufio"
	"encoding/base64"
	"encoding/csv"
	"errors"
	"flag"
	"os"

	"github.com/miekg/dns"
)

type FileEncodeCommand struct {
	fs *flag.FlagSet

	dnsType    string
	outputFile string
}

func NewFileEncodeCommand() *FileEncodeCommand {
	cmd := &FileEncodeCommand{
		fs: flag.NewFlagSet("file-encode", flag.ContinueOnError),
	}

	cmd.fs.StringVar(&cmd.dnsType, "t", "A", "Type of DNS request. (refer to types.go for supported types)")
	cmd.fs.StringVar(&cmd.outputFile, "o", "output.csv", "Output CSV file.")

	return cmd
}

func (cmd *FileEncodeCommand) Name() string {
	return cmd.fs.Name()
}

func (cmd *FileEncodeCommand) Init(args []string) error {
	return cmd.fs.Parse(args)
}

func (cmd *FileEncodeCommand) Run() error {
	args := cmd.fs.Args()
	if len(args) < 1 {
		return errors.New("missing filename")
	}

	filename := args[0]
	return cmd.encodeAndWrite(cmd.dnsType, filename)
}

func (cmd *FileEncodeCommand) encodeAndWrite(t string, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	s := bufio.NewScanner(f)

	tt, err := mapType(t)
	if err != nil {
		return err
	}

	r := make([][]string, 0, 100)

	for s.Scan() {
		hostname := s.Text()

		msg := dns.Msg{}
		msg.SetQuestion(normalizeHostname(hostname), tt)

		wire, err := msg.Pack()
		if err != nil {
			return err
		}

		r = append(r, []string{hostname, base64.StdEncoding.EncodeToString(wire)})
	}

	csvFile, err := os.Create(cmd.outputFile)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	err = csvWriter.Write([]string{"hostname", "encoded"})
	if err != nil {
		return err
	}

	err = csvWriter.WriteAll(r)
	if err != nil {
		return err
	}

	return nil
}
