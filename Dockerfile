FROM frolvlad/alpine-glibc

COPY chopper /usr/bin/chopper
ADD resources /resources
EXPOSE 8000
WORKDIR /
ENTRYPOINT ["chopper"]
