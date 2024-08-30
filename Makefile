build:
	mkdir -p output
	CCO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o output/wallpaper-bin_arm main.go
	CCO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o output/wallpaper-bin_amd main.go
	makefat ./output/wallpaper-bin ./output/wallpaper-bin_*
	rm -rf ./output/wallpaper-bin_*

clean:
	rm -rf ./output
	rm -rf ./wallpaper


run: build
	./output/wallpaper-bin