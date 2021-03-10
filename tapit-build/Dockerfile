FROM ubuntu

ENV PROVISION_CONTEXT "production"

# Deploy scripts/configurations
COPY static/ /static
COPY tapit /tapit
COPY entrypoint.sh /entrypoint.sh

run apt-get update
run apt-get install -y ca-certificates

# Harder to bypass
ENTRYPOINT ["/entrypoint.sh"]

# Can be overwritten -- with run args
CMD ["/entrypoint.sh"]
