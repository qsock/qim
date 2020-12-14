import os

from proj import confgo, logicgo, maingo, handlego, makefile


def gen(name, src_dir):
    gen_dir(name, src_dir)
    gen_gofile(name, src_dir)

def gen_dir(name, src_dir):
    srv_dir = src_dir + "/server/%s_server" % (name)
    conf_dir = srv_dir + "/config"
    logic_dir = srv_dir + "/logic"
    os.makedirs(conf_dir, 0o755, exist_ok=True)
    os.makedirs(logic_dir, 0o755,  exist_ok=True)

def gen_gofile(name, src_dir):
    srv_dir = src_dir + "/server/%s_server" % (name)
    conf_dir = srv_dir + "/config"
    logic_dir = srv_dir + "/logic"

    confgo.gen(name, conf_dir)
    logicgo.gen(name, logic_dir)
    maingo.gen(name, srv_dir)
    handlego.gen(name, srv_dir)
    makefile.gen(name, src_dir)