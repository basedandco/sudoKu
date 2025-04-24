FROM ubuntu:25.04

# Install necessary packages
RUN apt-get update && apt-get install -y \
    build-essential \
    libpam0g-dev \
    golang \
    sudo \
    neovim \
    pamtester \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Create test user
RUN useradd -m testuser && \
    echo "testuser:password" | chpasswd && \
    usermod -aG sudo testuser

# Set up working directory
WORKDIR /sudoku

# Copy project files
COPY . .

# Build and install the PAM module
RUN make build && \
    make install

# Add a test script to verify functionality
RUN echo '#!/bin/bash\necho "Try running: sudo -u testuser sudo ls"\necho "This should trigger the sudoku authentication"' > /test.sh && \
    chmod +x /test.sh

# Start with a shell
CMD ["/bin/bash"]
