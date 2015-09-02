FROM busybox:latest
MAINTAINER Julien Levesy <julien@upfluence.co>

ADD etcdexpose-bin /bin/etcdexpose
RUN chmod +x /bin/etcdexpose

ENTRYPOINT ["/bin/etcdexpose"]
