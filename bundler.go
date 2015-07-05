package bundler

import (
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"
)

// var err error

// ByFile bundles the file and writes the files back
// out on w, writer
func ByFile(filePath string, w io.Writer) error {

	input, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	// file, err := os.Open(filePath)
	// if err != nil {
	// 	return err
	// }
	// defer file.Close()

	return bundle(bytes.NewReader(input), filepath.Dir(filePath), w)
}

// ByIo bundles the file input by r, the reader
// and writes the files back out on w, writer
// assuming file bundle paths are relative to the
// provided path string
func ByIo(r io.Reader, dir string, w io.Writer) error {

	var err error

	if !filepath.IsAbs(dir) {
		if dir, err = filepath.Abs(dir); err != nil {
			return err
		}
	}

	return bundle(r, dir, w)
}

// Bundle bundles the file input by r, the reader
// and writes the files back out on w, writer
func bundle(r io.Reader, absPath string, w io.Writer) error {

	var input []byte
	var err error

	input, err = ioutil.ReadAll(r)
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
			path := absPath + "/" + itm.val

			input, err = ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			if err = bundle(bytes.NewReader(input), filepath.Dir(path), w); err != nil {
				return err
			}
		case itemEOF:
			break LOOP
		}
	}

	return nil
}
