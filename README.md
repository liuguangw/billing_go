# billing_go

这是一个用Go语言编写的billing验证服务器，可以处理用户登录、注册、退出和查询点数

> 开源版本的billing不提供点数兑换的源码

##### 支持以下任意一种运行环境:
- linux 32位(Linux 2.6.23 或者以上版本，并且已安装glibc)

- Windows 7, Server 2008R2 或者更高版本

##### 下载地址:

[https://github.com/liuguangw/billing_go/releases](https://github.com/liuguangw/billing_go/releases)

#### 编译方法
建议在Windows下编译，首先需要安装好go开发环境，然后安装mysql所需的go语言库

```bash
#安装go-sql-driver/mysql(需要已安装Git环境)
go get -u github.com/go-sql-driver/mysql
```

最后双击build.bat即可构建

```
billing       - linux版本的billing服务器
billing.exe   - windows版本的
config.json  - 配置文件
```

#### 部署方法

这里只讲一下如何部署linux版本的billing服务器

首先把 `billing` 和`config.json`上传到服务端的任意位置比如`/home`目录，然后修改`config.json`中的相关配置

```json
{
  "ip": "127.0.0.1",//billing服务器的ip，默认127.0.0.1即可
  "port": 12680,//billing服务器监听的端口(自定义一个未被占用的端口即可)
  "db_host": "127.0.0.1",//MySQL服务器的ip或者主机名
  "db_port": 3306,//MySQL服务器端口
  "db_user": "root",//MySQL用户名
  "db_password": "root",//MySQL密码
  "db_name": "web",//账号数据库名(一般为web)
  "allow_ips": []//允许的服务端连接ip,为空时表示允许任何ip,不为空时只允许指定的ip连接
}
```

`billing` 和`config.json`必须放置于同一目录下

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

```bash
#启动billing方法

#进入billing所在文件夹,比如/home
cd /home
#添加执行权限
chmod a+x ./billing
#启动billing，并让其在后台运行
./billing &
```

