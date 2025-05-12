# billing_go

这是一个用 Go 语言编写的 billing 验证服务器。

## 运行环境要求

支持以下任意一种运行环境

- Linux (Linux 内核版本 2.6.32 或者以上版本)

- Windows 10, Server 2016 或者更高版本

  > 参考 Go 语言程序的系统需求，https://go.dev/wiki/MinimumRequirements

## bug 反馈

如果使用此程序出现问题，可以提交 Issue 说明你所遇到的问题，并附上相关日志文件。

**billing** 的日志在 billing 程序所在的目录，文件名为 billing.log，首次运行时会自动创建此文件。

此外还需要 **Login** 服务器的日志，因为只有 Login 服务器会连接 billing 服务器。

可以修改一下运行脚本 `run.sh`，把 Login 服务器的日志写入某个文件，这样就可以在问题复现时从日志中查询到原因。

```sh
# 记录Login服务器日志的方法

# 原命令
./Login >/dev/null 2>&1 &

# 修改/dev/null 为 /home/login.log
# 修改后,日志文件会保存到/home/login.log
./Login >/home/login.log 2>&1 &
```

## 获取程序包

有以下两种方式获取此程序包。

### 1.使用我编译好的版本

[点击这里](https://github.com/liuguangw/billing_go/releases)查看我编译好的各版本, 可能没有最新的编译版本

### 2.手工编译

如果你想亲自进行编译，需要确保你的操作系统满足以下条件

- 设备已连接网络
- 已安装 Git
- 已安装 make(仅 linux 需要)
- 已安装 Go 1.23 或者更高版本

linux 使用 make 进行编译

```bash
# make命令说明

# 构建
make

# 清理
make clean

# 构建并且打包32位
make x32

# 构建并且打包64位
make x64

# 构建并且打包所有架构
make all
```

windows 下直接双击 build.bat 进行编译

## 相关文件说明

```
billing       - Linux版本的billing服务器
billing.exe   - Windows版本的
config.yaml   - 配置文件
```

## 配置文件

配置文件和程序必须放在同一个目录下，配置文件支持两种格式`yaml`或者`json`，配置文件名称为`config.yaml`或者`config.json`，如果两个文件都存在，则`yaml`格式优先。

### yaml 格式的配置示例

```yaml
# #后面的为注释
# 字符串可以不加引号,除非里面有#字符,所以如果数据库密码有#字符、空格时,就要加上引号
#
#billing服务器的ip，默认127.0.0.1即可
ip: 127.0.0.1
#
#billing服务器监听的端口(自定义一个未被占用的端口即可)
port: 12680
#
#MySQL服务器的ip或者主机名
db_host: localhost
#
#MySQL服务器端口
db_port: 3306
#
#MySQL用户名
db_user: root
#
#MySQL密码
db_password: 'root'
#
#账号数据库名(一般为web)
db_name: web
#
#只有在老版本MySQL报old_password错误时,才需要设置为true
allow_old_password: false
#
#用户登录的账号不存在时,是否引导用户进行注册
auto_reg: true
#
#允许的服务端连接ip,为空时表示允许任何ip,不为空时只允许指定的ip连接,
#allow_ips:
#  - 1.1.1.1
#  - 127.0.0.1
#
#点数修正, 在查询点数时,显示的值多1点或者少1点时才需要配置此值, 正常情况设置为0
#如果显示的点数少1点则配置为1(这是临时方案,一般是lua脚本有问题,lua脚本不应该把返回的数值减一,建议修改客户端查询点数的脚本)
#如果显示的点数多一点则配置为-1(一般不可能发生)
point_fix: 0
#登录的玩家总数量限制,如果为0则表示无上限
max_client_count: 500
#
#每台电脑最多可以同时登录的用户数量限制,如果为0则表示无上限
pc_max_client_count: 3
# billing类型 0经典 1怀旧
bill_type: 0
```

### json 格式的配置示例

```json
{
  "ip": "127.0.0.1",
  "port": 12680,
  "db_host": "localhost",
  "db_port": 3306,
  "db_user": "root",
  "db_password": "root",
  "db_name": "web",
  "allow_old_password": false,
  "auto_reg": true,
  "allow_ips": [],
  "point_fix": 0,
  "max_client_count": 500,
  "pc_max_client_count": 3,
  "bill_type": 0
}
```

> 如果 biiling 和服务端位于同一台服务器的情况下，建议 billing 的 IP 使用 127.0.0.1,这样可以避免绕一圈外网
>
> 本项目中附带的配置文件各项值为其默认值,如果你的配置中的字段的值与默认值相同,则可以省略相同的字段配置

将`billing` (Windows 服务器则是`billing.exe`)和配置文件放置于同一目录下

修改游戏服务器的配置文件`....../tlbb/Server/Config/ServerInfo.ini`中 billing 的配置

```ini
#........
[Billing]
Number=1
#billing服务器的ip
IP0=127.0.0.1
# billing服务器监听的端口
Port0=12680
#.........
```

最后启动游戏服务端、启动 billing 即可

## 启动与停止

Linux 和 Windows 下的操作分别如下

### 启动

Linux 下启动 billing 方法(**前台模式**)

```bash
#进入billing所在文件夹,比如/home
cd /home
#添加执行权限
chmod +x ./billing
#启动billing
./billing
```

Linux 以守护进程后台运行 billing 的方法(**daemon 模式**)

```bash
#进入billing所在文件夹,比如/home
cd /home
#添加执行权限
chmod +x ./billing
#启动billing
./billing up -d
```

Linux 以 **systemd** 方式运行 billing 服务的方法:
参考文件 [billing.service](billing.service)

Windows 下,直接双击`billing.exe`即可

### 停止

停止 billing 命令

```bash
# 使用stop命令
./billing stop

# 也可以使用kill命令
kill -SIGTERM $(pgrep -f billing)
```

如果是前台模式，可以使用 Ctrl + C 组合键停止服务器
