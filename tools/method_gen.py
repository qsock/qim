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

    for i in range (0, len(proto_list)):
        path = os.path.join(proto_dir, proto_list[i])
        with open(path, "r") as f:
            fname = get_file_name(f.name)
            flines = []
            for line in f:
                line = line.lstrip()
                if not line.startswith("rpc"):
                    continue
                pattern = re.compile(r'[a-zA-Z0-9]+')
                data = pattern.findall(line)
                if len(data) < 3:
                    continue
                # 取第二个作为methodname
                api_name = data[1]
                if api_name=="Ping" or len(api_name) == 0:
                    continue
                method_name = " "+fname.capitalize()+api_name
                method_name += " = "
                method_name += "\"/%s.%s/%s\"\n"%(fname,fname.capitalize(),api_name)
                flines.append(method_name)
            if len(flines)==0 :
                continue
            content = "package method\n\n"
            content += "const (\n"
            for fline in flines:
                content += fline
            content += ")\n"
            fpath = os.path.join(src_dir, "lib", "method",fname+".go")
            with open(fpath,"w") as f2:
                f2.write(content)
#             print(fpath,content)

def get_file_name(path_string):
    pattern = re.compile(r'([^<>/\\\|:""\*\?]+)\.\w+$')
    data = pattern.findall(path_string)
    if data:
        return data[0]

if __name__ == "__main__":
    main(sys.argv)