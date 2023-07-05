# Makefile:

# Generate Resources
bundle_svgs:
	 fyne bundle --package commons -o ./internal/commons/svgImages.go ./resources

package_mac:
	fyne package -os darwin -icon skoona.png --name m2i