package main

import (
	"net/http"
	"net/url"

	tokenCode "github.com/PBH-Tech/moonenv/lambdas/endpoints/auth"
	restApi "github.com/PBH-Tech/moonenv/lambdas/util/rest-api"
)

func RevokeToken(deviceCode string, refreshToken string) (restApi.Response, error) {
	token, err := tokenCode.GetToken(deviceCode)

	if err != nil {
		return restApi.BuildErrorResponse(http.StatusNotFound, "Device code not found")
	}

	return invalidateToken(token.ClientId, refreshToken)
}

func invalidateToken(clientId string, refreshToken string) (restApi.Response, error) {
	// TODO: I don't like this code that are very similar to refresh too, should we abstract it?
	data := url.Values{}
	oauthUrl, err := url.ParseRequestURI(CognitoUrl)

	if err != nil {
		return restApi.BuildErrorResponse(http.StatusInternalServerError, "Error while parsing cognito url")
	}

	data.Set("client_id", clientId)
	data.Set("token", refreshToken)
	oauthUrl.Path = "/oauth2/revoke"

	oauthUrlStr := oauthUrl.String()
	client := &http.Client{}

	resp, err := client.PostForm(oauthUrlStr, data)

	if err != nil {
		return restApi.BuildErrorResponse(http.StatusInternalServerError, "Error while sending HTTP request")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return restApi.BuildErrorResponse(http.StatusUnauthorized, "Invalid refresh token")
	}

	return restApi.ApiResponse(http.StatusNoContent, nil)
}
