FROM golang as builder
WORKDIR /go/src/oauth
RUN go get github.com/joho/godotenv
RUN go get github.com/satori/go.uuid
RUN go get github.com/streadway/amqp
RUN go get golang.org/x/crypto/bcrypt
RUN go get github.com/dgrijalva/jwt-go
RUN go get github.com/TonPC64/gomon
RUN go get github.com/globalsign/mgo
RUN go get github.com/go-chi/chi
RUN go get golang.org/x/oauth2
RUN go get cloud.google.com/go/compute/metadata
RUN go get github.com/thedevsaddam/govalidator