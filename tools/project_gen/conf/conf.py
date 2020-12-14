import os

toml = '''# 对外暴露的http端口号
port=8000
# 是什么环境
env="{env}"

# etcd的配置
[op]
    # etcd集群的地址
    endpoints      = ["127.0.0.1:2379"]
    # 注册的prefix
    prefix         = "/main/{env}"
    # 注册的服务名称
    server_name    = "{server_name}"
    # 监听的服务名称
    watch_servers  = []

# 日志类型，文件：file，终端：stdout，系统日志：syslog
logtype = "file"
# 日志等级
loglevel = 1
[log]
    # 文件名称
    filename       = "{file_name}"
    # 文件数量
    filenum        = 50
    # 文件大小Mb
    filesize       = 256
    # 日志的输出文件夹
    dir            = "/srv/{file_name}/logs"
    #  是否开启gzip压缩
    use_gzip       = true'''

def gen(name, src_dir):
    op_name = name+'.'+name.capitalize()
    server_name = name+"_server"

    dev_name = src_dir+"/config/dev/" + server_name+".toml"
    prod_name = src_dir + "/config/prod/" + server_name+".toml"

    dev_toml = toml.format(env='dev',
                           name=name,
                           server_name=op_name,
                           file_name=server_name)
    prod_toml = toml.format(env='prod',
                           name=name,
                           server_name=op_name,
                           file_name=server_name)

    with open(dev_name, 'w') as f:
        f.write(dev_toml)
    with open(prod_name, 'w') as f:
        f.write(prod_toml)