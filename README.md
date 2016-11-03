# Go Docker Router Proxy

License: [MIT](./LICENSE.md)

## System Requirements

Pathes an filenames in this guide are for **Mac OS** *(but with some modifications this will work on any *nix system)*

* Go *(tested with version 1.7)*
* dnsmasq *(recommended)*
* apache2 *(recommended)*

## Goal

If you are working with a lot of service-container, messing around with ports within your URLs might become unhandy.

So we created an litte script that will take care of this and allows you to handle HTTP/S-requests in an easier way.

**Note:** Please be aware, that SSL termination has to be handled by yourself.  

## How It Works

We will create an double reverse proxy, so we are able to resolve dynamic container configurations while still having an local webserver listening on port 80.

To achive this, we configure DNSmasq so it's capable of resolving an custom domain *(``.dock`` in our case but, you can choose an different one)* 
and tell our OS to use our local DNS server for this domain.

Then we configure our existing Apache2 server to listen to any ``.dock`` domain and to proxying this request to our ``docker-router-proxy`` tool.

Our tool will be called and tries to resolve the requested domain based on started container names.

## Command arguments

* ``-debug`` Will print some informations for debugging in case of issues
* ``-port`` Defines on wich port the proxy will listen, by default this will be ``:9800``
* ``-domain`` By defaut ``.dock`` will be used but with this parameter you can choose your own

## Examples

**Our current Docker process-list:**

```
CONTAINER ID   IMAGE          COMMAND       CREATED             STATUS             PORTS                  NAMES
28ddf8dba251   ubuntu:16.04   "/bin/bash"   About an hour ago   Up About an hour   0.0.0.0:8999->80/tcp   drunk_wescoff
```

**Our HTTP request:**

```
GET http://drunk-wescoff.dock/
```

**Internal name resolving:**

1. Normalizing: ``drunk-wescoff.dock`` becomes ``drunk-wescoff``
2. Lookup ``drunk-wescoff``: no match
3. Fallback lookup by replacing *all* dashes with underscores (``drunk_wescoff``): Match for ``0.0.0.0:8999``
4. Proxy the current request to ``0.0.0.0:8999``


## Install

**1. Make shure you configured your ``$GOBIN`` variable**

```
export GOBIN="/usr/local/bin"
```

**2. Run installation**

```
go install docker-router-proxy.go
```

**3. Starting docker-router-proxy**

```
docker-router-proxy
```

*Hint: If you need ``docker-router-proxy`` on a regular base, it's recommended to add this tool to your startup-items so it will be startet once you logged on.*

## DNSmasq Configuration

In this step we will create a virtual domain named ``dock`` so our Docker containers will be available with names like ``http://my-fancy-container.dock/``.

**1. Add this configuration to your ``dnsmasq.conf`` configuration**

```
address=/dock/127.0.0.1
```

**2.1 Make shure the ``/etc/resolver``-directory exists**

```
mkdir -p /etc/resolver
```

**2.2 Add an dedicated resolver for this new domain**

```
echo "nameserver 127.0.0.1" > /etc/resolver/dock
```

## Apache2 Configuration

**1. Enable ``mod_mod_proxy`` in your ``httpd.conf``**

```
LoadModule proxy_module libexec/mod_proxy.so
LoadModule proxy_connect_module libexec/mod_proxy_connect.so
LoadModule proxy_http_module libexec/mod_proxy_http.so
```

**2. Configure your reverse-proxy as vhost**

```
<VirtualHost *:80>
    ProxyPreserveHost On
    ProxyPass / http://0.0.0.0:9800/
    ProxyPassReverse / http://0.0.0.0:9800/

    ServerName my-container.dock
    ServerAlias *.dock
</VirtualHost>
```

**3. Restart your apache**

```
apachectl restart
```

## Known issues

*Issue:* At the moment there are some priority issues with containers that can serve HTTP and HTTPS traffic.

*Workaround:* Create an container for each protocol
