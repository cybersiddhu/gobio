package seqio

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"regexp"
)

type Fasta struct {
	Id       []byte
	Sequence []byte
}

type FastaReader struct {
	reader      *bufio.Reader
	fastaRegExp *regexp.Regexp
	seenHeader  bool
	header      []byte
	sequence    []byte
	entry       *Fasta
	exhausted   bool
}

func (f *FastaReader) NextEntry() *Fasta {
	return f.entry
}

func (f *FastaReader) HasEntry() bool {
	for {
		line, err := f.reader.ReadSlice('\n')
		if err == io.EOF {
			if !f.exhausted {
				f.exhausted = true
				f.entry = &Fasta{
					Id: f.header, Sequence: f.sequence,
				}
				return true
			}
			return false
		}
		if bytes.HasPrefix(line, []byte(">")) {
			if match := f.fastaRegExp.FindSubmatch(line); match != nil {
				if !f.seenHeader {
					f.header = match[1]
					f.seenHeader = true
				} else {
					f.entry = &Fasta{Id: f.header, Sequence: f.sequence}
					f.header = match[1]
					f.sequence = []byte{}
					return true
				}
			}
		} else {
			f.sequence = append(f.sequence, bytes.TrimSuffix(line, []byte("\n"))...)
		}
	}
	return false
}

func NewFastaReader(file string) *FastaReader {
	reader, err := os.Open(file)
	if err != nil {
		panic(err.Error())
	}
	return &FastaReader{
		reader:      bufio.NewReader(reader),
		fastaRegExp: regexp.MustCompile(`^>(\S+)`),
	}
}
