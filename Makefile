# Variables to compile go client
GO_CLIENT=bitwrk-client
GO_RELEASE_DIST=./dist
CLIENT_LINUX=$(GO_RELEASE_DIST)/bitwrk_linux_amd64/$(GO_CLIENT)
CLIENT_DARWIN=$(GO_RELEASE_DIST)/bitwrk_darwin_amd64/$(GO_CLIENT)
CLIENT_WINDOWS=$(GO_RELEASE_DIST)/bitwrk_windows_amd64/$(GO_CLIENT).exe
ADDON_NAME_ROOT=render_bitwrk
VERSION=0.7.0

# Variables to zip addon and client daemon into one 
TMPDIR=tmp
CLIENT_DIR=bitwrk_client
RESOURCE_DIR=resources
RENDER_DIR=bitwrk-blender
ADDON_DIR=render_bitwrk

all: build-go prep-addon package-darwin package-linux package-windows cleanup-addon

build-go:
	goreleaser --snapshot --skip-publish --rm-dist

prep-addon:
	echo "CLEAN UP PREVIOUS BUILD"
	rm -rf $(TMPDIR)
	echo "CREATE ADDON DIRECTORY STRUCTURE"
	mkdir $(TMPDIR) && \
	cp -r $(RENDER_DIR)/$(ADDON_DIR) $(TMPDIR)/ && \
	mkdir $(TMPDIR)/$(ADDON_DIR)/$(CLIENT_DIR) && \
	cp -r $(RESOURCE_DIR) $(TMPDIR)/$(ADDON_DIR)/$(CLIENT_DIR)/

package-darwin:
	echo "DARWIN: COPY CLIENT EXECUTABLE TO ADDON STRUCTURE"
	cp $(CLIENT_DARWIN) $(TMPDIR)/$(ADDON_DIR)/$(CLIENT_DIR)/ && \
	echo "DARWIN: ZIP ADDON" && \
	cd $(TMPDIR) && \
	zip -r ../$(ADDON_NAME_ROOT)_osx.zip * && \
	echo "DARWIN: REMOVE CLIENT EXECUTABLE" && \
	cd .. && \
	rm $(TMPDIR)/$(ADDON_DIR)/$(CLIENT_DIR)/$(GO_CLIENT)

package-linux:
	echo "LINUX: COPY CLIENT EXECUTABLE TO ADDON STRUCTURE"
	cp $(CLIENT_LINUX) $(TMPDIR)/$(ADDON_DIR)/$(CLIENT_DIR)/ && \
	echo "LINUX: ZIP ADDON" && \
	cd $(TMPDIR) && \
	zip -r ../$(ADDON_NAME_ROOT)_linux.zip * && \
	echo "LINUX: REMOVE CLIENT EXECUTABLE" && \
	cd .. && \
	rm $(TMPDIR)/$(ADDON_DIR)/$(CLIENT_DIR)/$(GO_CLIENT)

package-windows:
	echo "WINDOWS: COPY CLIENT EXECUTABLE TO ADDON STRUCTURE" && \
	cp $(CLIENT_WINDOWS) $(TMPDIR)/$(ADDON_DIR)/$(CLIENT_DIR)/ && \
	echo "WINDOWS: ZIP ADDON" && \
	cd $(TMPDIR) && \
	zip -r ../$(ADDON_NAME_ROOT)_windows.zip * && \
	echo "WINDOWS: REMOVE CLIENT EXECUTABLE" && \
	cd .. && \
	rm $(TMPDIR)/$(ADDON_DIR)/$(CLIENT_DIR)/$(GO_CLIENT).exe

cleanup-addon:
	rm -rf $(TMPDIR)
