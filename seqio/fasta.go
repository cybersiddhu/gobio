package seqio

import (
	"bufio"
	"bytes"
	"os"
	"regexp"
)

type FastaEntry struct {
	Id       []byte
	Sequence []byte
}

type FastaScanner struct {
	reader      *bufio.Reader
	fastaRegExp *regexp.Regexp
	header      []byte
	sequence    []byte
}

var seenHeader bool = false

func FastaSplitter(data []byte, atEOF bool) (advance int, token []byte , err error) {
	 advance, token, error = bufio.ScanLines(data,atEOF)
	 if err != nil {
	 		return
	 }

}

func (f *FastaScanner) NextSeq() (*Fasta, error) {
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

func NewFastaScanner(file string) *FastaScanner {
	reader, err := os.Open(file)
	if err != nil {
		panic(err.Error())
	}
	return &FastaReader{
		reader:      bufio.NewReader(reader),
		fastaRegExp: regexp.MustCompile(`^>(\S+)`),
	}
}
