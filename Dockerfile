FROM frolvlad/alpine-glibc

COPY squirrelchopper /usr/bin/squirrelchopper
ADD resources /resources
EXPOSE 8000
WORKDIR /
ENTRYPOINT ["squirrelchopper"]
