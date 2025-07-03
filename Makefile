dev:
	CGO_ENABLED=0 go build -ldflags="-s -w" -tags=debug

release:
	CGO_ENABLED=0 go build -ldflags="-s -w" -tags=release

clean:
	rm jte
	rm jte.log
