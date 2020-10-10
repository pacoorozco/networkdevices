# networkdevices


## Get all devices
```shell script
$ curl http://localhost:8010/devices
```

## Create a new device
```shell script
$ curl -d '{"fqdn":"my-hostname.domain.com.","model":"ios-xe", "version":"11.2DS"}' -H 'Content-Type: application/json' http://localhost:8010/devices
```

## Get a device
```shell script
$ curl http://localhost:8010/devices/my-hostname.domain.com.
```


## Update a device
```shell script
$ curl -d '{"fqdn":"my-hostname.domain.com.", "version":"11.5DS"}' -X PUT -H 'Content-Type: application/json' http://localhost:8010/devices
```

## Delete a device
```shell script
$ curl -X DELETE http://localhost:8010/devices/my-hostname.domain.com.
```
