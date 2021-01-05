---
title: "Basic Remote Setup"
date: 2020-02-27T20:35:00+01:00
weight: 16
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: basic-remote-setup.md
---

{{< toc >}}

The default configuration or the oCIS binary and the `owncloud/ocis` docker image is assuming, that you access oCIS on `localhost`. This enables you to do quick testing and development without any configuration.

If you need to access oCIS on a VM, docker container or a remote machine via an other hostname than `localhost`, you need to configure this hostname in oCIS. The same also applies if you are not using hostnames, but an IP instead (eg. `127.0.0.1`).

### Start the oCIS fullstack server

As a preparation you ne

In the following examples you have the binary in your current working directory, it is named 'ocis' and is marked as executable.


In order to run oCIS with self generated certificates please execute following command:
```bash
OCIS_LOG_LEVEL=WARN
KONNECTD_LOG_LEVEL=DEBUG
PROXY_HTTP_ADDR=0.0.0.0:443 \
OCIS_URL=https://ocis.owncloud.test:9200 \
sudo ./ocis server
```

When you have your own certificates in place, you also may running following command:
```bash
PROXY_HTTP_ADDR=0.0.0.0:9200 \
OCIS_URL=https://your-host:9200 \
PROXY_TRANSPORT_TLS_KEY=./certs/your-host.key \
PROXY_TRANSPORT_TLS_CERT=./certs/your-host.crt \
./bin/ocis server
```

{{< hint info >}}
**TLS Certificate**\
In this example, we are replacing the default self-signed cert with a CA signed one to avoid the certificate warning when accessing the login page.
{{< /hint >}}


For more configuration options check the configuration section in [ocis](https://owncloud.github.io/ocis/configuration/) and every ocis extension.



## Use Docker Compose

Please have a look at our other [deployment examples]({{< ref "./_index.md" >}}).
