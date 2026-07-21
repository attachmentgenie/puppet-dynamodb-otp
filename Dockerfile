FROM scratch
ARG TARGETPLATFORM
ENTRYPOINT ["/usr/bin/puppet-dynamodb-otp"]
COPY $TARGETPLATFORM/puppet-dynamodb-otp /usr/bin/
