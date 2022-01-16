FROM ubuntu
ENV PORT=80
ENV VERSION="v1.1"
LABEL multi.tech="http" multi.type="exercise" other="toy"
ADD bin/amd64/httpServer /httpServer
EXPOSE 80
ENTRYPOINT /httpServer