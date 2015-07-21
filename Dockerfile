FROM busybox:latest
MAINTAINER Julien Levesy <julien@upfluence.co>

ADD https://github.com/upfluence/etcdexpose/releases/download/v0.0.5/etcdexpose-linux-amd64-0.0.5 /bin/etcdexpose
RUN chmod +x /bin/etcdexpose

ENTRYPOINT ["/bin/etcdexpose"]
