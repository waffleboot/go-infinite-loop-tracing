
.PHONY: dlv parse trace all .FORCE

.FORCE:

all: main trace parse

main: main.go .FORCE
	@go build --gcflags='-l -N' -o main main.go

trace: main
	@GOMAXPROCS=1 ./main

dlv: main
	GOMAXPROCS=1 ~/go/bin/dlv exec loop

parser: parse.go parser.go order.go goroutines.go
	go build -o $@ $^

#traceGoPreempt
#goyield_m
#gopreempt_m

parse: parser
	@./parser > log.txt
