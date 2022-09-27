FROM golang:1.18-bullseye as baseExe
WORKDIR /backend
COPY . .
RUN go mod download
RUN go build -o /formfiller

FROM ubuntu:22.04
WORKDIR /app
RUN apt-get update && apt install -y libreoffice
RUN apt-get install -y python3-pip && pip3 install unoserver 
COPY ./resources /usr/share/fonts/truetype
COPY --from=baseExe /formfiller /formfiller
CMD ["/formfiller"]
