.PHONY: all install clean

TARGETNAMEWIN				:= old-jew-licenseserver-x64.exe
TARGETNAMELINUX				:= old-jew-licenseserver-x64-unix
TARGETNAMEARM				:= old-jew-licenseserver-arm64
TARGETNAMEMAC				:= old-jew-licenseserver-x64-darwin

all: windows linux arm64 darwin

windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o $(TARGETNAMEWIN) *.go

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o $(TARGETNAMELINUX) *.go

arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o $(TARGETNAMEARM) *.go

darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o $(TARGETNAMEMAC) *.go

clean:
	rm $(TARGETNAMELINUX) $(TARGETNAMEWIN) $(TARGETNAMEARM) $(TARGETNAMEMAC)

