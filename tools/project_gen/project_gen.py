#!/usr/bin/python3
import os
import sys

from conf import conf
from proj import proj
from proto import proto


def main(args):
    args = args[1:]
    if len(args) == 0:
        print("name should be assigned")
        exit(3)
    name = args[0]
    src_dir = "../../"
    os.chdir(src_dir)
    print(os.getcwd())
    src_dir = os.getcwd()

    if len(args) > 1:
        src_dir = args[1]
    name = name.lower()
    proto.gen(name, src_dir)
    proj.gen(name, src_dir)
    conf.gen(name, src_dir)

    os.system('make p;')


if __name__ == "__main__":
    main(sys.argv)