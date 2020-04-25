from invoke import task

@task
def update_protobug_service(c):
    with c.cd("./router"):
        c.run("protoc ./router.proto --go_out=plugins=grpc:./")

@task
def build(c):
    c.run("docker-compose build")
    c.run("docker-compose push")
