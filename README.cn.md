
# Httphere
一个开发者工具。  
在本地任何目录快速启动一个HTTP服务。  
本地多项目开发辅助工具。

## 作用

 - 任何目录启动web服务，默认输出index.html
 - 自动转发到后台服务，方便调整后端接口，同时避免跨域问题
 - 不同域名可指向不同服务器。
 - 指定后缀可指向不同服务器。


# 安装
查看 docs中 quick_start.md

### 使用场景1
分享当下目录文件供他人下载，或者浏览网页
```bash
httphere
```


### 使用场景2
当下目录为前台文件， 其他请求转发到后台服务  
接口类请求会转发到127.0.0.1:8080服务。
```bash
httphere -backend=127.0.0.1:8080
```

### 使用场景3
微服务场景下， 后台有多个服务， 根据httphere.yaml.example在当前目录创建httphere.yaml,  
配置hosts项， httphere启动时会自动加载当前目录下httphere.yaml配置文件。


### 使用场景4
手机端图片上传到电脑：
http://{{你的局域网IP}}/httphere_upload


