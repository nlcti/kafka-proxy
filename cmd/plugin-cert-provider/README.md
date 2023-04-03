# cert-provider plugin config

## use refreshable certificate via certificate provider

```
make clean build plugin.cert-provider

build/kafka-proxy server \
            --bootstrap-server-mapping "localhost:19092,0.0.0.0:30001" \
            --bootstrap-server-mapping "localhost:29092,0.0.0.0:30002" \
            --bootstrap-server-mapping "localhost:39092,0.0.0.0:30003" \
            --cert-provider-plugin-enable  \
            --cert-provider-plugin-command=build/cert-provider  \
            --cert-provider-plugin-param=--updated-proxy-listener-cert-file=/var/run/secrets/proxytls/tls.crt  \
            --cert-provider-plugin-param=--updated-proxy-listener-key-file=/var/run/secrets/proxytls/tls.key  \
            --cert-provider-plugin-param=--update-check-interval-minutes=25
```

The parameters `--updated-proxy-listener-cert-file` and `--updated-proxy-listener-key-file` specify paths to the required files of a x509-key pair. This key pair may be updated periodically.
The parameter `--update-check-interval-minutes` therefore gives control over the frequency of checks for certificate renewal.

All parameters are mandatory.
