FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ADD usermgr /
ENTRYPOINT ["/usermgr"]
CMD ["web"]

