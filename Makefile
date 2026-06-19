IMAGE       ?= node:lts
REPO        ?= $(CURDIR)
SANDBOX_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

.PHONY: sandbox

sandbox:
	@command -v smolvm >/dev/null 2>&1 || \
		{ echo "smolvm not found. Install: curl -sSL https://smolmachines.com/install.sh | bash"; exit 1; }
	smolvm machine run --net -it \
		--image $(IMAGE) \
		-v "$(REPO):/workspace" \
		-v "$(SANDBOX_DIR)scripts:/sandbox" \
		-- /bin/sh /sandbox/entrypoint.sh
