# OpenConfig coverage

A simple example to undertstand OpenConfig coverage of your IOS XR router configs.

![coverage](static/oc-diff.svg)

## Requirements to run

### Enviroment variables

- Target device details:

- USER (E.g. root)
- PASSWORD (E.g. mypassword)
- ADDRESS (E.g. [2001:420:2cff:1204::5501:3]:57344)

```console
export USER=cisco
export PASSWORD=cisco
export ADDRESS='[2001:420:2cff:1204::5501:3]:57344'
```

### Router config

The device needs to have gRPC enabled.

```console
grpc
 port 57344
 address-family ipv6
```

### Router's TLS Certificate

In order to secure the connection, a certificate file will nedd to be retrieved form the device.

```console
RP/0/RP0/CPU0:mrstn-5501-3.cisco.com# bash cat /misc/config/grpc/ems.pem
Tue Jul  3 14:43:05.712 UTC
-----BEGIN CERTIFICATE-----
MIIEEDCCAvigAwIBAgICFJgwDQYJKoZIhvcNAQENBQAwgbAxCzAJBgNVBAYTAlVT
MQswCQYDVQQIEwJDQTERMA8GA1UEBxMIU2FuIEpvc2UxFzAVBgNVBAkTDjM3MDAg
...
-----END CERTIFICATE-----
```