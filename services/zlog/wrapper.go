package zlog

import "io"

type writerWrapper struct {
	io.Writer
}

func (wrapper *writerWrapper) Sync() error {
	return nil
}
