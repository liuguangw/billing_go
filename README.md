# billing_go[dev分支]

这是一个用Go语言编写的billing验证服务器，此分支正在开发中，可能会有未知的问题。

## 运行环境要求

支持以下任意一种运行环境
- Linux 32位(Linux 2.6.23 或者以上版本)
- Windows 7, Server 2008R2 或者更高版本

## 获取程序包

  有以下两种方式获取此程序包。

  ### 1.使用我编译好的版本

  [点击这里](https://github.com/liuguangw/billing_go/releases)查看我编译好的各版本

  ### 2.手工编译
  如果你想亲自进行编译，需要确保你的操作系统满足以下条件

  - 设备已连接网络

  - 已安装Git

  - 已安装Go 1.17或者更高版本

    #### Windows环境下编译

    在Windows下编译，只需要双击项目目录下的`build.bat`即可
    
    #### Linux环境下编译

    ```bash
    # 切换到项目目录下执行下面的命令
    go build -o billing
    ```

## 相关文件说明

```
billing       - Linux版本的billing服务器
billing.exe   - Windows版本的
config.yaml   - 配置文件
```

## 配置文件

配置文件和程序必须放在同一个目录下，配置文件支持两种格式`yaml`或者`json`，配置文件名称为`config.yaml`或者`config.json`，如果两个文件都存在，则`yaml`格式优先。

### yaml格式的配置示例

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
#兑换参数，有的版本可能要设置为1才能正常兑换,有的则是1000
transfer_number: 1000
#
#登录的玩家总数量限制,如果为0则表示无上限
max_client_count: 500
#
#每台电脑最多可以同时登录的用户数量限制,如果为0则表示无上限
pc_max_client_count: 3
```



### json格式的配置示例

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
  "transfer_number": 1000,
  "max_client_count": 500,
  "pc_max_client_count": 3
}
```

> 如果biiling和服务端位于同一台服务器的情况下，建议billing的IP使用127.0.0.1,这样可以避免绕一圈外网
>
> 本项目中附带的配置文件各项值为其默认值,如果你的配置中的字段的值与默认值相同,则可以省略相同的字段配置
>

将`billing` (Windows服务器则是`billing.exe`)和配置文件放置于同一目录下

修改游戏服务器的配置文件`....../tlbb/Server/Config/ServerInfo.ini`中billing的配置

```ini
#........
[World]
IP=127.0.0.1
Port=777

[Billing]
Number=1
#billing服务器的ip
IP0=127.0.0.1
# billing服务器监听的端口
Port0=12680
#.........
```

最后启动游戏服务端、启动billing即可

## 启动与停止
Linux和Windows下的操作分别如下

### 启动

Linux下启动billing方法(**前台模式**)

```bash
#进入billing所在文件夹,比如/home
cd /home
#添加执行权限
chmod a+x ./billing
#启动billing
./billing
```

Linux下在后台运行billing的方法(**后台模式**,只支持Linux和类unix系统)
```bash
#进入billing所在文件夹,比如/home
cd /home
#添加执行权限
chmod a+x ./billing
#启动billing
./billing up -d
```

Windows下,直接双击`billing.exe`即可

### 停止

停止billing命令

```bash
./billing stop
```

如果是前台模式，可以使用Ctrl + C 组合键停止服务器
