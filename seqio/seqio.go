package seqio

import (
	"bufio"
	"bytes"
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
}

func (f *FastaReader) NextSeq() (*Fasta, error) {
	for {
		line, err := f.reader.ReadSlice('\n')
		if err != nil {
			if !f.seenHeader {
				f.seenHeader = true
				return &Fasta{
					Id: f.header, Sequence: f.sequence,
				}, nil
			}
			return nil, err
		}
		if match := f.fastaRegExp.FindSubmatch(line); match != nil {
			if !f.seenHeader {
				f.header = match[1]
				f.seenHeader = true
			} else {
				fasta := &Fasta{Id: f.header, Sequence: f.sequence}
				f.header = match[1]
				f.seenHeader = false
				return fasta, nil
			}
		} else {
			f.sequence = append(f.sequence, bytes.TrimSuffix(line, []byte("\n"))...)
		}
	}
	return nil, nil
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
