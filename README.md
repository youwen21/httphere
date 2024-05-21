# Httphere

Run a http server at any directory

## Quickstart
 - download binary executable [Download Now](https://github.com/youwen21/httphere)
 - move httphere to /usr/local/bin

Grant execution permission
```bash
chmod +x httphere
```

## Support platform：
 - Linux
 - Windows
 - MacOS

## Functions
 - A server view files
 - A http server serve static pages
 - Upload file to server root directory
 - A reverse proxy

### More detail about Upload file
 After startup ，terminal print link this:
```bash
view url: http://192.168.3.5:8089
upload url: http://192.168.3.5:8089/httphere_upload
qr url: http://192.168.3.5:8089/qr

```
Other computer or phone in LAN network can upload files to server root directory


## Command arguments
support arguments list
- host  (default: 0.0.0.0)
- port   (default: 80)
- backend   

```bash
httphere -port=8000 -backend="http://xxx.com"
```
Request forward to backend if the resource not find in server root directory


# Advanced config

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
