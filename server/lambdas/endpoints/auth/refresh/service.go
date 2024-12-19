package main

import (
	"encoding/json"
	"net/http"
	"net/url"

	tokenCode "github.com/PBH-Tech/moonenv/lambdas/endpoints/auth"
	restApi "github.com/PBH-Tech/moonenv/lambdas/util/rest-api"
)

type RefreshTokenBodyRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type CognitoRefreshTokenResponse struct {
	IdToken     string `json:"id_token"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type APIResponse struct {
	IdToken     string `json:"idToken"`
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
	TokenType   string `json:"tokenType"`
}

func RefreshToken(deviceCode string, request restApi.Request) (restApi.Response, error) {
	var data RefreshTokenBodyRequest

	err := json.Unmarshal([]byte(request.Body), &data)

	if err != nil {
		return restApi.BuildErrorResponse(http.StatusBadRequest, "Invalid body request. Refresh token field is required")
	}

	token, err := tokenCode.GetToken(deviceCode)

	if err != nil {
		return restApi.BuildErrorResponse(http.StatusNotFound, "Device code not found")
	}

	return getToken(token.ClientId, data.RefreshToken)
}

func getToken(clientId string, refreshToken string) (restApi.Response, error) {
	data := url.Values{}
	oauthUrl, err := url.ParseRequestURI(CognitoUrl)

	if err != nil {
		return restApi.BuildErrorResponse(http.StatusInternalServerError, "Error while parsing cognito url")
	}

	data.Set("client_id", clientId)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	oauthUrl.Path = "/oauth2/token"

	oauthUrlStr := oauthUrl.String()
	client := &http.Client{}

	resp, err := client.PostForm(oauthUrlStr, data)

	if err != nil {
		return restApi.BuildErrorResponse(http.StatusInternalServerError, "Error while sending HTTP request")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return restApi.BuildErrorResponse(http.StatusInternalServerError, "HTTP request did not return OK")
	}

	var tokenResponse CognitoRefreshTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return restApi.BuildErrorResponse(http.StatusInternalServerError, "Error decoding JSON response")
	}

	return restApi.ApiResponse(http.StatusCreated, APIResponse(tokenResponse))
}
