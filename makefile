default:
	@tmux_send "make build"

build:
	@go build main.go
	@./main
	@img_view "out.png"
