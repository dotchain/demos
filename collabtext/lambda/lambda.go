// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dotchain/dot/x/nw"
	"strings"
)

func main() {
	handler := &nw.Handler{Store: nw.MemPoller(nw.MemStore(nil))}
	lambda.Start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		var resp events.APIGatewayProxyResponse
		ct := request.Headers["Content-Type"]
		httpError := func(status string, code int) {
			resp.StatusCode = code
		}
		body := &strings.Builder{}
		handler.HandleLambda(ctx, ct, httpError, strings.NewReader(request.Body), body)
		resp.Body = body.String()
		return resp, nil
	})
}
