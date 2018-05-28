# irpcx
irpcx base on rpcx

# Usage

## Server

[irpcx-server](https://github.com/youngbloood/irpcx/tree/master/_example/server)

also you can define youself server just like:
```
    import("github.com/smallnets/rpcx/server")
    serve := server.New(basePath, []string{"127.0.0.1:2379"})
    serve.Serve=server.Server
```
    

## Client
[irpcx-client](https://github.com/youngbloood/irpcx/tree/master/_example/client)
