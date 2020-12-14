
content = '''package logic'''

def gen(name, srv_dir) :
    with open(srv_dir+"/logic.go","w") as f:
        f.write(content)