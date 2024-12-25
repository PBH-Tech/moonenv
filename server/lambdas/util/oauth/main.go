package oauth

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	restApi "github.com/PBH-Tech/moonenv/lambdas/util/rest-api"
)

var (
	cognitoUrl = os.Getenv("CognitoUrl")
)

func GetOAuthUrl() (*url.URL, *restApi.Response) {
	println()
	oauthUrl, err := url.ParseRequestURI(fmt.Sprintf("https://%s", cognitoUrl))

	if err != nil {
		response := restApi.BuildErrorResponse(http.StatusInternalServerError, "Error while parsing cognito url")
		return nil, &response
	}
	return oauthUrl, nil
}
