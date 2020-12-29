#!/usr/bin/env make
.PHONY : all release clean

release:
ifndef TAG
	@echo "Please provide a tag with TAG=v0.0.42"
	exit 1
endif
	sed -i '' 's|github.com/tophergopher/mongotest.*|github.com/tophergopher/mongotest $(TAG)|g' go.mod
	git add go.mod
	git commit -m "Updating revision to $(TAG)."
	git tag $(TAG)
	git push --tags
	cd ../mongotest && sed -i '' 's|github.com/tophergopher/easymongo.*|github.com/tophergopher/easymongo $(TAG)|g' go.mod
	cd ../mongotest && git add go.mod && git commit -m "Updating revision to $(TAG)."
	cd ../mongotest && git tag $(TAG) && git push --tags && git push --all
	git push --all