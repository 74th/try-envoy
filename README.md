# Envoy のテスト

## 0. 作業用コンテナ

```
kubectl apply -k manifests/0_working
kubectl exec -it po/working -- bash
```

## 1. gRPC のサービスのデプロイ

```
kubectl apply -k manifests/1_gRPC
```

ポートフォワードで動作することを確認

```
kubectl port-forward deploy/service 8080
```

Kubernetes 内で service の名前で DNS を引くと、2 つの Pod が返る

```
# nslookup server
Server:		10.152.183.10
Address:	10.152.183.10#53

Name:	server.envoy.svc.cluster.local
Address: 10.1.6.17
Name:	server.envoy.svc.cluster.local
Address: 10.1.6.19
```

## 2. Envoy Proxy のデプロイ

[./envoy.yaml](./manifests/2_envoy/envoy.yaml)

stern で見ると、1 つの gRPC でラウンドロビンされていることが確認できた。

## 3. 偏った係数でのラウンドロビン

[./envoy.yaml](./manifests/3_use_load_balancing_weight/envoy.yaml)

server1 と server2 の 2 つのサービスを作る。
clusters の所に、`load_balancing_weight`を書くだけ。

```yaml
clusters:
  - name: server
    connect_timeout: 0.5s
    type: STRICT_DNS
    dns_lookup_family: V4_ONLY
    lb_policy: ROUND_ROBIN
    http2_protocol_options: {}
    load_assignment:
      cluster_name: server
      endpoints:
        - lb_endpoints:
            - endpoint:
                address:
                  socket_address:
                    address: server1
                    port_value: 8080
              load_balancing_weight: 1
            - endpoint:
                address:
                  socket_address:
                    address: server2
                    port_value: 8080
              load_balancing_weight: 5
```

[./envoy.yaml](./manifests/2_envoy/envoy.yaml)
