package bundler

import (
	"bytes"
	"io"
	"io/ioutil"
)

// Bundle bundles the file input by r, the reader
// and writes the files back out on w, writer
func Bundle(r io.Reader, w io.Writer) error {

	input, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	l := lex("bundler", string(input), "", "")

LOOP:
	for {
		itm := l.nextItem()

		switch itm.typ {
		case itemText:
			w.Write([]byte(itm.val))
		case itemFile:
			input, err := ioutil.ReadFile(itm.val)
			if err != nil {
				return err
			}

			Bundle(bytes.NewReader(input), w)
		case itemEOF:
			break LOOP
		}
	}

	return nil
}
