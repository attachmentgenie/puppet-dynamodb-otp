FROM scratch
ENTRYPOINT ["/puppet-dynamodb-otp"]
COPY puppet-dynamodb-otp /