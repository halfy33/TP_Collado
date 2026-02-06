# Image golang pour la compilation
FROM golang:1.25 AS builder

# Copier les fichiers dont ont a besoin
COPY src /go/src
COPY Makefile /go

# Compilation
RUN make build

# Image finale
FROM busybox

# Installer le logiciel compilé
COPY src/www /www
COPY src/.env /
COPY src/servers.yaml /
COPY --from=builder /go/bin/buildup /buildup

# Partages
EXPOSE 80

# Démarrage
ENTRYPOINT [ "/buildup" ]