#!/usr/bin/python3
import os
import sys
import re


def main(args):
    args = args[1:]
    if len(args) == 0:
        print("name should be assigned")
        exit(3)
    src_dir = args[0]
    proto_dir = os.path.join(src_dir, "proto")
    if not os.path.exists(proto_dir):
        print(proto_dir+" is not our path")
        exit(3)
    proto_list = os.listdir(proto_dir)
    fcontent = "package controller\nconst (\n"
    for i in range (0, len(proto_list)):
        path = os.path.join(proto_dir, proto_list[i])
        with open(path, "r") as f:
            fname = get_file_name(f.name)
            fcontent += fname+"Proto = `"+f.read()+"`\n"
    fcontent += "\n)"
    with open(os.path.join(src_dir,"api_gateway","controller","proto.go"),"w") as f :
        f.write(fcontent)

def get_file_name(path_string):
    pattern = re.compile(r'([^<>/\\\|:""\*\?]+)\.\w+$')
    data = pattern.findall(path_string)
    if data:
        return data[0]

if __name__ == "__main__":
    main(sys.argv)