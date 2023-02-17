
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


# 本地开发使用httphere有哪些好处

### 不用修改nginx配置
有些情况，预览一个网站要增加nginx配置。 使用httphere避免增加配置的麻烦， 小白用户使用更方便。

### 配置代理更方便
只需要支持backend，后端http://ip:port， 不存在请求便转发给后端。

### php开发者不用配置nginx
支持FastCGI协议， httphere可以直接通过fastCGI与php-fpm交互。

### 调试更方便，
前端把便宜好的disc发给后端，不用配置proxy，不用考跨域问题。

### 共享方件更方便
在要共享的文件夹启动httphere,  局域网用户可以通过http服务下载对应的文件。