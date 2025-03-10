# Stage 1 : Compilation du programme Go
FROM golang:latest AS builder

# Définir le répertoire de travail dans le conteneur de compilation
WORKDIR /app

# Copier les fichiers sources du programme Go
COPY go.mod main.go ./
#COPY . .

# Télécharger les dépendances Go
RUN go mod download

# Compiler le programme Go
RUN CGO_ENABLED=0 GOOS=linux go build -v -o myapp

# Stage 2 : Création de l'image finale
#FROM alpine:latest
FROM scratch

# Copier le binaire compilé depuis le stage de compilation
COPY --from=builder /app/myapp /myapp

# Copier le dossier de données
COPY js /js
COPY css /css

# Définir le répertoire de travail dans l'image finale
WORKDIR /

# Exécuter le programme Go au démarrage du conteneur
CMD ["/myapp"]

