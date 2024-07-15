FROM scratch
ENTRYPOINT ["/puppet-dynamodb-otpe"]
COPY puppet-dynamodb-otpe /