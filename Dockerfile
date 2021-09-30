FROM ubuntu:20.04

ADD https://artifacts.elastic.co/downloads/logstash/logstash-7.10.2-amd64.deb /tmp/logstash-7.10.2-amd64.deb
RUN dpkg -i /tmp/logstash-7.10.2-amd64.deb && \
    rm -f /tmp/logstash-7.10.2-amd64.deb

COPY logstash-pipeline-tester /usr/local/bin
COPY config.yaml /etc/logstash-pipeline-tester/config.yaml

EXPOSE 8080

CMD ["logstash-pipeline-tester"]