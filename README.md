vpntool
=======

## Install
```
go get github.com/trusch/vpntool
```

## Usage

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






