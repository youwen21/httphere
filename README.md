# Httphere

一个开发者工具。  
使用httphere方便于手机端文件上传到电脑， 电脑文件下载到手机端， 局域网同事传输文件等场景。  

支持：

- 浏览文件
- 局域网内上传，下载文件
- 静态网站服务器
- 反向代理 - 不存在的请求，作为代理服务器转发到后端
- php服务器，fast-cgi支持php脚本

# 安装

see docs/quick_start.md

# 启动http服务

```bash
httphere
```

## 配置
httphere会尝试从当前目录读取httphere.yaml配置。  
如果存在使用配置覆盖默认配置项。

### 基础静态文件服务器配置
查看文件，启动静态网站场景，使用基础静态文件服务器配置
```yaml
base:
  #  listen_host: "0.0.0.0"  # default=0.0.0.0
  listen_port: "80"  # 指定监听端口
  static_server: "open"  # open | close,  是否开启静态文件浏览功能
  static_root: "./" # 指定项目根目录
  dump_request: "yes"  # 输出请求信息

tls:
  cert_file: "" # ssl配置cert
  key_file: ""  # ssl配置key
```

### 指定后台代理接口
作为前端开发人员，





## 命令行参数

support command params

- host
- port
- backend

