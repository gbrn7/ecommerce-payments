FROM golang:1.23

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod tidy

COPY . .

COPY .env .

RUN go build -o ecommerce-payments

RUN chmod +x ecommerce-payments

EXPOSE 9003

CMD [ "./ecommerce-payments" ]