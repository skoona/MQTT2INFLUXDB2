# Makefile:

# Generate Resources
bundle_svgs:
	 fyne bundle --package commons -o ./internal/commons/svgImages.go ./resources

package_mac_gui:
	fyne package --tags gui -os darwin -icon skoona.png --name m2i

package_mac_cli:
	go build --tags cli -o bin/m2i_cli cmd/cli/main.go