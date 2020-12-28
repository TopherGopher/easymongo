#!/usr/bin/env make
.PHONY : all release clean

release:
ifndef TAG
	@echo "Please provide a tag with --tag"
	exit 1
endif
	git tag $(TAG) && git push --tags
	cd ../mongotest && git tag $(TAG) && git push --tags && git push --all
	git push --all