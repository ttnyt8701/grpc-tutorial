コンパイル  
```
protoc -I. --go_out=. --go-grpc_out=. proto/*.proto
```

証明書  
```
brew install mkcert  
mkcert -install  
mkcert -CAROOT
mkdir ssl  
cd ssl
mkcert localhost
```

