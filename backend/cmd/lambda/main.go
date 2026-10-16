package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/njsatish/dreamslayerz-portal/backend/internal/api"
	appconfig "github.com/njsatish/dreamslayerz-portal/backend/internal/config"
	"net/http"
	"net/http/httptest"
)

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	req, _ := http.NewRequestWithContext(ctx, request.RequestContext.HTTP.Method, request.RawPath, nil)
	req.URL.RawQuery = request.RawQueryString
	recorder := httptest.NewRecorder()
	cfg, err := appconfig.Load("config")
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500, Body: `{"error":"invalid configuration"}`}, nil
	}
	api.New(cfg).ServeHTTP(recorder, req)
	headers := map[string]string{}
	for key, values := range recorder.Header() {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}
	return events.APIGatewayV2HTTPResponse{StatusCode: recorder.Code, Headers: headers, Body: recorder.Body.String()}, nil
}
func main() { lambda.Start(handler) }
