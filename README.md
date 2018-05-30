## irpcx
irpcx base on [rpcx](https://github.com/smallnest/rpcx)

## Usage

### Server

[irpcx-server](https://github.com/youngbloood/irpcx/tree/master/_example/server)

also you can define youself server just like:
```
    import("github.com/smallnets/rpcx/server")
    serve := server.New(basePath, []string{"127.0.0.1:2379"})
    serve.Serve=server.Server
```
    

### Client
[irpcx-client](https://github.com/youngbloood/irpcx/tree/master/_example/client)
also you can define youself client just like:
```
    basePath := "/irpcx/youngbloood"
    service := "Test"
    method := "Add"
    
    req := irpcx.NewRequest(basePath, service, method)
    client.InitEtcdAddr([]string{"127.0.0.1:2379"})
    // sync invoke rpc
    resp, err := client.Call(req)
    // async invoke rpc
    resp, err =client.Go(req)
```

### Contribute
    fork and pull request