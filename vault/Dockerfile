FROM hashicorp/vault:1.16

COPY docker-entrypoint.sh /vault/docker-entrypoint.sh

RUN chmod +x /vault/docker-entrypoint.sh

ENTRYPOINT ["/vault/docker-entrypoint.sh"]
