
# Httphere
一个开发者工具。  
在本地任何目录快速启动一个HTTP服务


# 安装

## 使用go安装
go install gitee.com/youwen21/httphere

## 下载可执行文件
> https://gitee.com/youwen21/httphere/releases

提供windows, linux, mocOS 三个平台可执行文件。

## 编译安装
```bash 
git clone http://gitee.com/youwen21/httphere
cd httphere 
go mod tidy
go install
```

# 启动HTTP服务
在任意目录执行命令  httphere  

启动动输出效果。  
```text
port is 8080
backend URL is http://127.0.0.1:8098/
root is /Users/owen/html/test
Listening on [::]:8080

```

# 配置监听端口，后台代理转发等信息

配置参数支持三种方式
 - 命令行参数
 - export env, 如 export PORT=8089;
 - .env 配置文件

## 启动参数
```bash 
httphere  -port=8058
```

touch .env file

```env
cd {pwd}
touch .env

PORT=3100
BACKEND="http://127.0.0.1:8099/"
```

# 参数列表
## 项目参数
 - port  默认值：8090
 - root 默认值：当前目录
 - backend   默认值：无

## FastCGI参数
 - fastcgi :  false|true
 - fastcgi_proto : tcp | unix
 - fastcgi_addr :  ip:port | unix file 
 - fastcgi_root :  php 项目的目录位置


