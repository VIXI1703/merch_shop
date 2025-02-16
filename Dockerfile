FROM ubuntu:latest
LABEL authors="hvalo"

ENTRYPOINT ["top", "-b"]