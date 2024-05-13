# gobackup
> 简单的数据库**全量**定时备份系统, 目前还处于开发阶段, 建议使用于**开发测试环境**

## 入门开始

## Linux安装
### 条件
>  安装 mysqldump/pg_dump/mongodump
### 配置说明
config.yaml
```yaml
level: "debug"      # 日志等级
store: "/app/data"  # 存储
cron: "* * * * *"   # 定时
backup:
- type: mysql       # 数据库类型 (pg, mysql)
  host: 192.168.1.1 # 数据库地址
  port: 3306        # 数据库端口
  user: root        # 数据库用户
  password: root    # 数据库密码
- type: pg
  host: 192.168.1.1
  port: 5432
  user: root
  password: root
- type: mongo
  host: 192.168.1.1
  port: 27017
  user: ""
  password: ""
- type: scp
  user: root 
  password: root
  host: 192.168.52.147
  port: 22
  fromPath: /tmp/scp
```

### 命令行启动
```shell
./gobackup
```

## Windows安装
### 配置说明
config.yaml
```yaml
level: "debug"      # 日志等级
store: "C:\\Users\\admin\\Desktop\\goproject\\gobackup\\data" # 存储
cron: "* * * * *"   # 定时
backup:
- type: mysql       # 数据库类型 (pg, mysql)
  host: 192.168.1.1 # 数据库地址
  port: 3306        # 数据库端口
  user: root        # 数据库用户
  password: root    # 数据库密码
- type: pg
  host: 192.168.1.1
  port: 5432
  user: root
  password: root
- type: mongo
  host: 192.168.1.1
  port: 27017
  user: ""
  password: ""
- type: scp
  user: root 
  password: root
  host: 192.168.52.147
  port: 22
  fromPath: /tmp/scp
```
### 启动
```shell
./gobackup.exe      # 执行定时任务
./gobackup.exe -run # 执行一次
```


## 说明
### 支持数据库类型
###### mysql
###### postgres
###### mongo
###### scp

### 支持存储
###### local
