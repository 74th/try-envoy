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
