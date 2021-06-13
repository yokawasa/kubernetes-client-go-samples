.PHONY: clean all svclist nodelist podlist podlist-in-svc istio-get-vs istio-update-vs

.DEFAULT_GOAL := all

TARGETS=\
	svclist \
	nodelist \
	podlist \
	podlist-in-svc \
	istio-get-vs \
	istio-update-vs

cur := $(shell pwd)
docker_image_repo := kubernetes-client-go-samples
docker_image_tag := latest

${TARGETS}:
	cd ${cur}/$@ && GO111MODULE=on go build -o ${cur}/dist/$@ main.go

all: $(TARGETS)

clean:
	rm -rf dist

lint:
	golint ./...

docker-build:
	set -x
	export DOCKER_BUILDKIT=1
	docker build -t ${docker_image_repo}:${docker_image_tag} . --target executor

# docker-run: docker-build
docker-run:
	docker run --rm -it -v ${HOME}/.kube/config:/root/.kube/config ${docker_image_repo}:${docker_image_tag} sh -c "/svclist"
	docker run --rm -it -v ${HOME}/.kube/config:/root/.kube/config ${docker_image_repo}:${docker_image_tag} sh -c "/nodelist"
	docker run --rm -it -v ${HOME}/.kube/config:/root/.kube/config ${docker_image_repo}:${docker_image_tag} sh -c "/podlist"
	docker run --rm -it -v ${HOME}/.kube/config:/root/.kube/config ${docker_image_repo}:${docker_image_tag} sh -c "/podlist-in-svc"
	docker run --rm -it -v ${HOME}/.kube/config:/root/.kube/config ${docker_image_repo}:${docker_image_tag} sh -c "/istio-get-vs -s hoge-api -n hogens"
	docker run --rm -it -v ${HOME}/.kube/config:/root/.kube/config ${docker_image_repo}:${docker_image_tag} sh -c "/istio-update-vs -s hoge-api -n hogens -desthost hoge-api.hogens.svc.cluster.local -destsubset hoge-api"

# docker-run-aws: docker-build
docker-run-aws:
	docker run --rm -it -v ${HOME}/.kube/config:/root/.kube/config -v ${HOME}/.aws:/root/.aws ${docker_image_repo}:${docker_image_tag} sh -c "/svclist"
	docker run --rm -it -v ${HOME}/.kube/config:/root/.kube/config -v ${HOME}/.aws:/root/.aws ${docker_image_repo}:${docker_image_tag} sh -c "/nodelist"
	docker run --rm -it -v ${HOME}/.kube/config:/root/.kube/config -v ${HOME}/.aws:/root/.aws ${docker_image_repo}:${docker_image_tag} sh -c "/podlist"
	docker run --rm -it -v ${HOME}/.kube/config:/root/.kube/config -v ${HOME}/.aws:/root/.aws ${docker_image_repo}:${docker_image_tag} sh -c "/podlist-in-svc"
	docker run --rm -it -v ${HOME}/.kube/config:/root/.kube/config -v ${HOME}/.aws:/root/.aws ${docker_image_repo}:${docker_image_tag} sh -c "/istio-get-vs -s hoge-api -n hogens"
	docker run --rm -it -v ${HOME}/.kube/config:/root/.kube/config -v ${HOME}/.aws:/root/.aws ${docker_image_repo}:${docker_image_tag} sh -c "/istio-update-vs -s hoge-api -n hogens -desthost hoge-api.hogens.svc.cluster.local -destsubset hoge-api"
