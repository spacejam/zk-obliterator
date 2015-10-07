FROM debian
ADD ./zk-obliterator /work/
WORKDIR /work/
ENV ZK=master.mesos:2181
ENV SIZE=1024
ENV RATIO=0.8
ENV CONCURRENCY=40
CMD ./zk-obliterator -zk=$ZK -ratio=$RATIO -size=$SIZE -concurrency=$CONCURRENCY; sleep 10
