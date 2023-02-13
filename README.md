
# Httphere
a static file server, 
if file not exists, the request will handle by ReverseProxy.  
you can config the ReverseProxy backend url.

# install

## use go install
go install github.com/youwen21/httphere

## download execution directly 
> https://github.com/youwen21/httphere/releases


# configure
touch .env file 

```env
cd {pwd}
touch .env

PORT=3100
BACKEND="http://127.0.0.1:8099/"
```


# start server

in terminal
```bash
httphere
```

output
```text
port is 8080
backend URL is http://127.0.0.1:8098/
root is /Users/owen/html/test
Listening on [::]:8080

```

