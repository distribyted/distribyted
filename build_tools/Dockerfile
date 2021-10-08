FROM techknowlogick/xgo:go-1.17.x

# add 32-bit and 64-bit architectures and install 7zip
RUN \
    dpkg --add-architecture i386 && \
    dpkg --add-architecture amd64 && \
    apt-get update && \
    apt-get install -y --no-install-recommends p7zip-full

# install LIBFUSE
RUN \
    apt-get update && \
    apt-get install -y --no-install-recommends libfuse-dev:i386 && \
    apt-get install -y --no-install-recommends libfuse-dev:amd64 && \
    apt-get download libfuse-dev:i386 && \
    dpkg -x libfuse-dev*i386*.deb /

ENV \
    OSXCROSS_NO_INCLUDE_PATH_WARNINGS 1
