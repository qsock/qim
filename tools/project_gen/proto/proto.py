content='''syntax = "proto3";

package %s;
import "ret.proto";

option go_package="github.com/qsock/qim/lib/proto/%s";

service %s {
  rpc Ping (ret.NoArgs) returns (ret.NoArgs) {}
}
'''



def gen (name, src_dir):
    gen_file(name, src_dir)
    gen_makefile(name, src_dir)

def gen_file(name, src_dir):
    pkg = name.capitalize()
    str = content % (name, name, pkg)
    fname = src_dir + '/proto/' + name + '.proto'
    f = open(fname, 'w')
    f.write(str)
    f.close()


def gen_makefile(name, src_dir):
    # 追加
    mkfile = src_dir + '/Makefile'
    f = open(mkfile, 'r')
    lines = f.readlines()
    f.close()
    idx = 0
    for line in lines:
        if 'cp' in line:
            break
        idx += 1

    lines.insert(idx, '\tprotoc -I ./proto --gogofaster_out=plugins=grpc:. %s.proto;\n' % name)
    s = ''.join(lines)
    with open(mkfile, 'w') as f:
        f.write(s)
