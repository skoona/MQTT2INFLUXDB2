# Makefile:

# Generate Resources
bundle_svgs:
	 fyne bundle --package commons -o ./internal/commons/svgImages.go ./resources
