content = '''
%s:
	go install ${BUILDARG} github.com/qsock/qim/server/%s;\n'''

def gen(name, srv_dir) :

    server_name = name+"_server"
    contes = content % (server_name, server_name)
    with open(srv_dir+"/Makefile", "a") as f:
        f.write(contes)