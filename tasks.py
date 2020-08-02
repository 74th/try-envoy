from invoke import task

@task
def update_protobug_service(c):
    with c.cd("./router"):
        c.run("protoc ./router.proto --go_out=plugins=grpc:./")

@task
def build(c):
    c.run("docker-compose build")
    c.run("docker-compose push")

@task
def create_envoy_cert(c):
    with c.cd("manifests/4_gke_ingress/server1_cert"):
        c.run("go run /usr/local/go/src/crypto/tls/generate_cert.go -duration 876000h0m0s -host server1")
    with c.cd("manifests/4_gke_ingress/envoy_cert"):
        c.run("go run /usr/local/go/src/crypto/tls/generate_cert.go -duration 876000h0m0s -host envoy")
    with c.cd("router/server/cert"):
        c.run("go run /usr/local/go/src/crypto/tls/generate_cert.go -duration 876000h0m0s -host localhost")
