FROM mcr.microsoft.com/devcontainers/base:bookworm

COPY ZscalerRootCertificate.crt /usr/share/ca-certificates
RUN echo ZscalerRootCertificate.crt >> /etc/ca-certificates.conf
RUN update-ca-certificates
