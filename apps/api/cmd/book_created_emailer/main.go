package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/example/bookapi/internal/lambda/bookemailer"
)

func main() {
	handler, err := bookemailer.New(context.Background())
	if err != nil {
		log.Fatalf("configure handler: %v", err)
	}

	lambda.Start(handler.Handle)
}
