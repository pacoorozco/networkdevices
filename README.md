# networkdevices

![Test and build](https://github.com/pacoorozco/networkdevices/workflows/Test%20and%20build/badge.svg)

## Introduction
This is a sample application implementing, partially, [this CRUD API](https://app.swaggerhub.com/apis/pacoorozco/NetworkDevices/0.0.1) for a network device manager.

For more details you can read [this design document](https://docs.google.com/presentation/d/1fG3xQrDU_HUMct_D4YIQxmhlq3QrICNiwDLB-4i_rYI/edit?usp=sharing). 

## Getting started

1. Compile from the source.

You need golang installed in your system, and then:

```shell script
$ git clone https://github.com/pacoorozco/networkdevices.git
$ cd networkdevices
$ make build
```

2. Run the server
```shell script
$ ./deviceManagerAPI
```

3. Everything is ready to send requests to `http://localhost:8010`.

## Test the API

You can find the API definition [here](https://app.swaggerhub.com/apis/pacoorozco/NetworkDevices/0.0.1). Bear in mind that it's partially implemented.

> It's usind `FQDN` as device id to be used directly by network engineers. Other endpoints can be implemented using other id, such `UUID`.

### Get all devices
```shell script
$ curl http://localhost:8010/devices
```

### Create a device
```shell script
$ curl -d '{"fqdn":"my-hostname.domain.com.","model":"ios-xe", "version":"11.2DS"}' -H 'Content-Type: application/json' http://localhost:8010/devices
```

### Get a device
```shell script
$ curl http://localhost:8010/devices/my-hostname.domain.com.
```

### Update a device
```shell script
$ curl -d '{"fqdn":"my-hostname.domain.com.", "version":"11.5DS"}' -X PUT -H 'Content-Type: application/json' http://localhost:8010/devices
```

### Delete a device
```shell script
$ curl -X DELETE http://localhost:8010/devices/my-hostname.domain.com.
```
