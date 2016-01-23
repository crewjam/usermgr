
.PHONY: _usermgr lint

GO_SOURCES=$(shell find . -name \*.go -not -name index.go)

SOURCES=\
  $(GO_SOURCES) \
  web/index.html

GENERATED_SOURCES=\
	web/index.go

PLATFORM_BINARIES=\
	usermgr.Linux.x86_64 \
	usermgr.Linux.armv7l

IMAGE_NAME=crewjam/usermgr
GITHUB_USER=crewjam
GITHUB_REPOSITORY=usermgr

all: usermgr $(PLATFORM_BINARIES)

clean:
	-rm $(GENERATED_SOURCES)
	-rm usermgr $(PLATFORM_BINARIES)

web/index.go: web/index_gen.go web/index.html
	go generate ./...

usermgr:
	go build -a -installsuffix cgo -ldflags '-s' -o $@ ./cmd/usermgr/usermgr.go

usermgr.Linux.x86_64: $(SOURCES) $(GENERATED_SOURCES)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s' \
	  -o $@ ./cmd/usermgr/usermgr.go

usermgr.Linux.armv7l: $(SOURCES) $(GENERATED_SOURCES)
	GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s' \
	  -o $@ ./cmd/usermgr/usermgr.go

container: usermgr.linux.amd64
	docker build -t $(IMAGE_NAME) .

check:
	go test ./...

lint:
	go fmt ./...
	goimports -w $(GO_SOURCES)

release: lint check $(BINARIES) container
	@[ ! -z "$(VERSION)" ] || (echo "you must specify the VERSION"; false)
	which ghr >/dev/null || go get github.com/tcnksm/ghr
	#ghr -u $(GITHUB_USER) -r $(GITHUB_REPOSITORY) --delete $(VERSION) usermgr.linux.amd64
	#docker tag -f $(IMAGE_NAME) $(IMAGE_NAME):$(VERSION)
	#docker push $(IMAGE_NAME)
	#docker push $(IMAGE_NAME):$(VERSION)

foo:
	make usermgr.Linux.x86_64
	scp usermgr.Linux.x86_64 192.168.59.104:/tmp/usrmgr
	ssh 192.168.59.104 sudo cp /tmp/usermgr /opt/usermgr/bin/usermgr

appdeploy: $(GENERATED_SOURCES)
	~/project/go_appengine/goapp deploy -application usermgr-998 web/appengine/app.yaml
