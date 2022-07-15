FROM golang:1.18-bullseye as baseExe
WORKDIR /backend
# TODO: switch to multistage building for prod
RUN apt-get update && apt install -y libreoffice
RUN apt-get install -y python3-pip && pip3 install unoserver 
COPY . .
RUN go mod download
RUN go build -o /formfiller
CMD ["/formfiller"]
