# Binary make targets {
# Define variables
BINDIR := bin
SRCDIR := cmd
PKGDIR := pkg

# Find all Go source files in the cmd directory
SOURCES := $(wildcard $(SRCDIR)/*)

# Find all Go source files in the pkg directory
PACKAGES := $(wildcard $(PKGDIR)/*)

# Define the executables to be created
EXECUTABLES := $(patsubst $(SRCDIR)/%, $(BINDIR)/%, $(SOURCES))

# Default target
all: $(EXECUTABLES)

# Rule to build each executable
$(BINDIR)/%: $(SRCDIR)/% $(PACKAGES)
	@mkdir -p $(BINDIR)
	go build -o $@ ./$<

# Clean up binaries
clean:
	rm -rf $(BINDIR)

# PHONY targets
.PHONY: all clean

# Binary make targets }


# Container make targets {
CONTAINER_ENGINE := podman
IMAGE_NAME ?= authentication-server
REGISTRY ?=
REPOSITORY ?= localhost
TAG ?= latest


HOST_PORT := 8008
CONTAINER_PORT :=8008

VOLUME_NAME := authentication-server-data
VOLUME_MOUNT := /mnt/

# Build IMAGE_TAG dynamically based on which variables are set
IMAGE_TAG := $(if $(REGISTRY),$(REGISTRY)/)$(if $(REPOSITORY),$(REPOSITORY)/)$(IMAGE_NAME)$(if $(TAG),:$(TAG))

.PHONY: build-image clean-containers container-volume clean-volume run-container-foreground

build-image:
	$(CONTAINER_ENGINE) build -t $(IMAGE_TAG) .

clean-containers:
	$(CONTAINER_ENGINE) system prune -f

volume:
	$(CONTAINER_ENGINE) volume create $(VOLUME_NAME)
	tar -cf assets.tar assets/
	$(CONTAINER_ENGINE) volume import $(VOLUME_NAME) assets.tar
	rm assets.tar

clean-volume:
	$(CONTAINER_ENGINE) volume rm $(VOLUME_NAME)

run-container-foreground:
	$(CONTAINER_ENGINE) run --rm -p $(HOST_PORT):$(CONTAINER_PORT) -v $(VOLUME_NAME):$(VOLUME_MOUNT) $(IMAGE_TAG)

# Container make targets }
