package main

import (
	"net"

	"github.com/Microsoft/go-winio"
)

const (
	pipePath   = `\\.\pipe\foo`
	socketPath = "/tmp/foo.socket"
)

var connTypes map[string]connType = map[string]connType{
	"npipe": {
		listen: func() (net.Listener, error) {
			return winio.ListenPipe(pipePath, nil)
		},
		dial: func() (net.Conn, error) {
			return winio.DialPipe(pipePath, nil)
		},
	},
	"unix": {
		listen: func() (net.Listener, error) {
			return net.Listen("unix", socketPath)
		},
		dial: func() (net.Conn, error) {
			return net.Dial("unix", socketPath)
		},
	},
}
