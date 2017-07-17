package fakes

import (
	"io"

	"golang.org/x/exp/io/i2c/driver"
)

// FakeOpener is a fake i2c/driver.Opener with buffers that we can inspect.
type FakeOpener struct {
	W io.WriteCloser
	R io.ReadCloser
}

func (f *FakeOpener) Open(addr int, tenbit bool) (driver.Conn, error) {
	return &FakeConn{f.W, f.R}, nil
}

// FakeConn is a fake i2c/driver.Conn with buffers that we can inspect.
type FakeConn struct {
	W io.WriteCloser
	R io.ReadCloser
}

func (f *FakeConn) Tx(w, r []byte) error {
	var err error
	if w != nil {
		_, err = f.W.Write(w)
		if err != nil {
			return err
		}
	}
	if r != nil {
		_, err = f.R.Read(r)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *FakeConn) Close() error {
	if err := f.W.Close(); err != nil {
		return err
	}
	return f.R.Close()
}
