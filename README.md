# Go Wake-On-Lan in http server
What is Wake-on-lan? 
https://en.wikipedia.org/wiki/Wake-on-LAN

Wake up your computer on the LAN

# Usage
## Place the `devices.json` files in the same directory as the binaries
Change `devices.json` to appropriate address
```json
{
  "devices": [
    { "device_name":"Device1", "mac_address": "AA:BB:CC:DD:EE:FF", "broad_cast_ip": "192.168.1.255" },
    { "device_name":"Device1", "mac_address": "AA:BB:CC:DD:EE:FF", "broad_cast_ip": "192.168.1.255" }
  ]
}
```
## Build script
### For current machine
```
go build
```
### For mt7620 mipsel 
 powershell command
- specific os and CPU arch 

```shell script
$env:GOOS="linux"
$env:GOARCH="mipsle"
$env:GOMIPS="softfloat"
```

- static linking
```shell script
$env:CGO_ENABLED=0
```
- build
```
go build -o mipsel_wol  -trimpath -ldflags "-s -w -buildid=" .
```

### flag info 
- https://johng.cn/cgo-enabled-affect-go-static-compile/
- https://zhangsnow.com/go/golang%E7%BC%96%E8%AF%91%E5%8F%82%E6%95%B0ldflags.html
- https://segmentfault.com/a/1190000008323048
- `-s` disable symbol table 
- `-w` disable DWARF generation 
- `-trimpath`: remove absolute path in the trace





## Run script
### 
```
wol <port>
```
where port is http listen port, default is 9999

### example
```shell script
./wol 9998
```


# Thanks to
- https://github.com/sabhiram/go-wol.git
