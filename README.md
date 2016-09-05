vpntool
=======

## Install
```
go get github.com/trusch/vpntool
```

## Usage
```
Usage of vpntool:
  -clients string
    	add client(s) to vpn (accepts comma separated list)
  -deploy string
    	deploy this entity to --url
  -init
    	init vpn and create server
  -out string
    	ovpn directory (default ".")
  -pki string
    	pki directory (default "pki")
  -url string
    	url to use
```

### Initialize
```
vpntool --init
```

### Add Clients
```
vpntool --clients clientA,clientB,clientC --url my-vpn-server.com
```

### Deploy Server
```
vpntool --deploy server --url user@my-vpn-server.com
```

### Deploy Client
```
vpntool --deploy clientA --url user@clientA.com
```

### One-Shot-Setup
```
vpntool --init --clients clientA,clientB,clientC --url my-vpn-server.com
```
