## Makefile for sudoKu PAM module
# Define variables
NAME := pam_sudoku.so
SRC := sudoku.go
GO := go
GOOPTS := build -buildmode=c-shared -x -v -ldflags=-w
INCLUDE_DIR := /usr/include/security

# Detect system architecture and lib paths
UNAME_M := $(shell uname -m)
ifeq ($(UNAME_M),x86_64)
    # Check if /usr/lib/x86_64-linux-gnu exists first
    ifneq ($(wildcard /usr/lib/x86_64-linux-gnu),)
        LIBDIR := /usr/lib/x86_64-linux-gnu/security
    else ifneq ($(wildcard /lib/x86_64-linux-gnu),)
        LIBDIR := /lib/x86_64-linux-gnu/security
    else
        LIBDIR := /lib64/security
    endif
else ifeq ($(UNAME_M),i386)
    # Check if /usr/lib/i386-linux-gnu exists first
    ifneq ($(wildcard /usr/lib/i386-linux-gnu),)
        LIBDIR := /usr/lib/i386-linux-gnu/security
    else ifneq ($(wildcard /lib/i386-linux-gnu),)
        LIBDIR := /lib/i386-linux-gnu/security
    else
        LIBDIR := /lib/security
    endif
else
    LIBDIR := /usr/lib/security
endif

# Main targets
.PHONY: all build install test clean uninstall disable

all: build install

build: $(NAME)

$(NAME): $(SRC)
	$(GO) $(GOOPTS) -o $(NAME) $(SRC)

install: build
	@echo "Installing PAM module to $(LIBDIR)..."
	install -D -m 0644 $(NAME) $(LIBDIR)/$(NAME)
	@echo "Installing header file to $(INCLUDE_DIR)..."
	install -D -m 0644 pam_sudoku.h $(INCLUDE_DIR)/pam_sudoku.h
	@echo "Installation complete"

enable:
	@echo "Enabling PAM module..."
	@if ! grep -q "pam_sudoku.so" /etc/pam.d/sudo; then \
		sed -i '1s/^/auth    required    pam_sudoku.so\n/' /etc/pam.d/sudo; \
	fi	
	@echo "Module enabled"
	
disable:
	@echo "Disabling PAM module..."
	@if grep -q "pam_sudoku.so" /etc/pam.d/sudo; then \
		sed -i '/pam_sudoku.so/d' /etc/pam.d/sudo; \
	fi
	@echo "Module disabled"

uninstall: disable
	@echo "Uninstalling PAM module..."
	rm -f $(LIBDIR)/$(NAME)
	rm -f $(INCLUDE_DIR)/pam_sudoku.h
	@echo "Uninstallation complete"

clean:
	rm -f $(NAME)
	rm -f *.h

test: build
	@echo "Running test in Docker container..."
	docker build -t sudoku-pam-test -f Dockerfile .
	@echo "Starting Docker container with shell access..."
	@echo "Use 'exit' to leave the container when finished testing"
	docker run -it --rm sudoku-pam-test
