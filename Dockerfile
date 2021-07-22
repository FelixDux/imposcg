FROM golang:1.16.3-buster
ENV GO111MODULE=on
ENV GOFLAGS=-mod=mod
ENV APP_USER=app
ENV APP_HOME=/go/src/imposcg
ENV GIN_MODE=release
ENV PORT=8080
RUN mkdir -p $APP_HOME 
WORKDIR $APP_HOME
COPY . .
RUN go build
EXPOSE $PORT
ENTRYPOINT ["./imposcg"]
# CMD ["bee", "run"]
# CMD ["go", "run", "impact.oscillator.go"]