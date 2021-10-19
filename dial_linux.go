package main

import (
	"net"
)

const socketPath = "/tmp/foo.socket"

var connTypes map[string]connType = map[string]connType{
	"unix": {
		listen: func() (net.Listener, error) {
			return net.Listen("unix", socketPath)
		},
		dial: func() (net.Conn, error) {
			return net.Dial("unix", socketPath)
		},
	},
}
