
# Httphere
一个开发者工具。  
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
support command params
- host 
- port 
- backend

```bash
httphere -port=3000 -backend=127.0.0.1:80
```


# advance configure
cp httphere.yaml.example httphere.yaml
