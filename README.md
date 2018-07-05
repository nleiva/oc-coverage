# OpenConfig coverage

A simple example to understand OpenConfig coverage of your IOS XR router config. It will retrieve the original config of a device, then pull the OpenConfig version of it. The target device config will be replaced by a minimum setup to allow remote access, and then the OpenConfig config will be applied on top of it. The result will be compared to the initial state to determine the coverage gap.

![coverage](static/oc-diff.svg)

## Requirements to run

### Enviroment variables

Target device details:

- USER (E.g. root)
- PASSWORD (E.g. mypassword)
- ADDRESS (E.g. [2001:420:2cff:1204::5501:3]:57344)

You can setup this as follows:

```console
export USER=root
export PASSWORD=mypassword
export ADDRESS='[2001:420:2cff:1204::5501:3]:57344'
```

### gRPC config

The target device needs to have gRPC enabled.

```console
grpc
 port 57344
 address-family ipv6
```

### Router's base config

This is the config that will provide minimun settings to allow a remote gRPC connection; IP addressing, username/password, and gRPC setup. See the example at [base.cfg](input/base.cfg).

### Router's TLS Certificate

In order to secure the connection, a certificate file will need to be retrieved from the device.

```console
RP/0/RP0/CPU0:mrstn-5501-3.cisco.com# bash cat /misc/config/grpc/ems.pem
Tue Jul  3 14:43:05.712 UTC
-----BEGIN CERTIFICATE-----
MIIEEDCCAvigAwIBAgICFJgwDQYJKoZIhvcNAQENBQAwgbAxCzAJBgNVBAYTAlVT
MQswCQYDVQQIEwJDQTERMA8GA1UEBxMIU2FuIEpvc2UxFzAVBgNVBAkTDjM3MDAg
...
-----END CERTIFICATE-----
```
