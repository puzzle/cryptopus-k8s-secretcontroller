FROM scratch
USER nobody
COPY ./cryptopus-k8s-secretcontroller ./
ENTRYPOINT ["./cryptopus-k8s-secretcontroller", "-logtostderr"]
