package main

import (
	"bufio"
	"encoding/base64"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/miekg/dns"
)

type FileEncodeCommand struct {
	fs *flag.FlagSet

	outputFile string
}

func NewFileEncodeCommand() *FileEncodeCommand {
	cmd := &FileEncodeCommand{
		fs: flag.NewFlagSet("file-encode", flag.ContinueOnError),
	}

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

	return cmd.encodeAndWrite(args[0])
}

func (cmd *FileEncodeCommand) encodeAndWrite(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	s := bufio.NewScanner(f)

	r := make([][]string, 0, 100)

	for s.Scan() {
		data := strings.Split(s.Text(), ",")
		hostname := data[0]
		reqType := data[1]

		tt, err := mapType(reqType)
		if err != nil {
			fmt.Printf("unsupported type: %v\n", reqType)
			return err
		}

		msg := dns.Msg{}
		msg.SetQuestion(normalizeHostname(hostname), tt)

		wire, err := msg.Pack()
		if err != nil {
			return err
		}

		r = append(r, []string{hostname, reqType, base64.StdEncoding.EncodeToString(wire)})
	}

	csvFile, err := os.Create(cmd.outputFile)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	err = csvWriter.Write([]string{"hostname", "type", "encoded"})
	if err != nil {
		return err
	}

	err = csvWriter.WriteAll(r)
	if err != nil {
		return err
	}

	return nil
}
