
# Httphere
a static file server, 
if file not exists, the request will handle by ReverseProxy.  
you can config the ReverseProxy backend url.

# install
see docs/quick_start.md

# start server
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
