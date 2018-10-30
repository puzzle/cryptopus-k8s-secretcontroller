FROM scratch
USER nobody
COPY ./cryptopussecretcontroller ./
ENTRYPOINT ["./cryptopussecretcontroller", "-logtostderr"]
