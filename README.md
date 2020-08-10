# billing_go

这是一个用Go语言编写的billing验证服务器。

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

  - 已安装Go 1.12或者更高版本

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
config.json  - 配置文件
```

## 部署方法

修改`config.json`中的相关配置

```json
{
  "ip": "127.0.0.1",//billing服务器的ip，默认127.0.0.1即可
  "port": 12680,//billing服务器监听的端口(自定义一个未被占用的端口即可)
  "db_host": "127.0.0.1",//MySQL服务器的ip或者主机名
  "db_port": 3306,//MySQL服务器端口
  "db_user": "root",//MySQL用户名
  "db_password": "root",//MySQL密码
  "db_name": "web",//账号数据库名(一般为web)
  "allow_old_password": false,//只有在老版本MySQL报old_password错误时,才需要设置为true
  "auto_reg": true,//用户登录的账号不存在时,是否引导用户进行注册
  "allow_ips": [],//允许的服务端连接ip,为空时表示允许任何ip,不为空时只允许指定的ip连接
  "transfer_number": 1000 //兑换参数，有的版本可能要设置为1才能正常兑换,有的则是1000
}
```

> 如果biiling和服务端位于同一台服务器的情况下，建议billing的IP使用127.0.0.1,这样可以避免绕一圈外网
>
> 本项目中附带的`config.json`的各项值为其默认值,如果你的配置中的值与默认值相同,则可以省略
>
> 例如你的配置只有密码和端口和上方配置不同，则可以这样写
>
> {
>
>   "port" : 12681,
>
>   "db_password" : "123456"
>
> }
>
> 如果你的配置和默认配置完全一样，则可以简写为 {}

将`billing` (Windows服务器则是`billing.exe`)和`config.json`放置于同一目录下

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

Linux停止billing命令

```bash
./billing stop
```

Windows下关闭billing窗口即可

## 其它语言的实现

除了此版本外，我还有一个rust版本的实现，[点击这里](https://github.com/liuguangw/billing_rust)可以查看。