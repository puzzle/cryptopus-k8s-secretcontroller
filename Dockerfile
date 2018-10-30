FROM scratch
COPY ./vaultsecretcontroller ./
ENTRYPOINT ["./vaultsecretcontroller", "-logtostderr"]
