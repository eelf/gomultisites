## gomultisites ##
A simple go script for https terminating and reverse proxying for hosting different services on a single host

### binary ###
```
go install
ln -svf $(go env GOPATH)/bin/gomultisites /usr/local/sbin/
```

### config ###
```
go get github.com/eelf/jsonschemacheck
jsonschemacheck $(pwd)/config.schema.json $(pwd)/config.json
ln -svf $(pwd)/config.json /usr/local/etc/gomultisites.json
```

### rc ###
```
ln -svf $(pwd)/gomultisites /usr/local/etc/rc.d/
echo gomultisites_enable=yes >> /etc/rc.conf
```

### package ###
tbd
