
content = '''package main

import (
	"context"
	"github.com/qsock/qim/lib/proto/ret"
)

type Server struct{}

func (*Server) Ping(ctx context.Context, req *ret.NoArgs) (*ret.NoArgs, error) {
	return new(ret.NoArgs), nil
}'''

def gen(name, srv_dir) :
    with open(srv_dir+"/handle.go", "w") as f:
        f.write(content)