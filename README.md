
## binary ##
```sh
go install
ln -svf $(go env GOPATH)/bin/gomultisites /usr/local/sbin/
```

## config ##
```sh
go run jsonschemacheck/main.go $(pwd)/config.schema.json $(pwd)/config.json
ln -svf $(pwd)/config.json /usr/local/etc/gomultisites.json
```

## rc ##
```sh
ln -svf $(pwd)/gomultisites /usr/local/etc/rc.d/
```

## package ##
tbd

