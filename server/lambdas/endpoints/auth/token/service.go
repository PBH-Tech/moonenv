package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	tokenCode "github.com/PBH-Tech/moonenv/lambdas/endpoints/auth"
	restApi "github.com/PBH-Tech/moonenv/lambdas/util/rest-api"
	"github.com/google/uuid"
)

type CodeChallenge struct {
	CodeChallenge string
	CodeVerifier  string
}

type CognitoTokenResponse struct {
	IdToken      string `json:"id_token"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type APIResponse struct {
	IdToken      string `json:"idToken"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
	TokenType    string `json:"tokenType"`
}

func RequestSetOfToken(clientId string) (restApi.Response, error) {
	var (
		stateCode  = uuid.New().String()
		deviceCode = uuid.New().String()
		expiresIn  = 900 // 15 minutes
	)

	codeChallenge := generateCodeVerifierAndChallenge()
	authorizationUri := fmt.Sprintf(
		"%s/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&code_challenge=%s&code_challenge_method=S256&state=%s&scope=openid profile",
		CognitoUrl, clientId, CallbackUri, codeChallenge.CodeChallenge, stateCode)

	token, err := tokenCode.InsertToken(tokenCode.TokenCode{
		DeviceCode:              deviceCode,
		AuthorizationUri:        authorizationUri,
		VerificationUriComplete: fmt.Sprintf("%s/device?code=%s&authorize=true", CallbackUri, stateCode),
		ClientId:                clientId,
		CodeChallenge:           codeChallenge.CodeChallenge,
		CodeVerifier:            codeChallenge.CodeVerifier,
		Status:                  "authorization_pending",
		State:                   stateCode,
		ExpireAt:                strconv.FormatInt(time.Now().Add(time.Duration(expiresIn)*time.Second).Unix(), 10),
		LastCheckedAt:           strconv.FormatInt(time.Now().Unix(), 10),
	})

	if err != nil {
		return restApi.ApiResponse(http.StatusBadRequest, err.Error())
	}

	//TODO: don't return the whole token object
	return restApi.ApiResponse(http.StatusCreated, token)
}

func RequestJWTs(deviceCode string, clientId string) (restApi.Response, error) {
	token, err := tokenCode.GetToken(deviceCode)

	if token == nil || err != nil {
		return restApi.BuildErrorResponse(http.StatusNotFound, "Device code was not found")
	}

	response, err := validateTokenCode(*token, clientId)

	if response != nil || err != nil {
		return *response, err
	}

	err = tokenCode.UpdateToken(token.DeviceCode, tokenCode.TokenCode{LastCheckedAt: strconv.FormatInt(time.Now().Unix(), 10)})

	if err != nil {
		return restApi.BuildErrorResponse(http.StatusInternalServerError, "Failed to update token")

	}

	if token.Status == "authorization_pending" {
		return restApi.BuildErrorResponse(http.StatusBadRequest, "Authorization is still pending")
	} else if token.Status == "completed" {
		return restApi.BuildErrorResponse(http.StatusBadRequest, "This authorization flow is already completed")
	} else if token.LoginCode == "" {
		return restApi.BuildErrorResponse(http.StatusInternalServerError, "Something went wrong while setting the login code")
	} else {

		err = tokenCode.UpdateToken(token.DeviceCode, tokenCode.TokenCode{Status: "completed"})

		if err != nil {
			return restApi.BuildErrorResponse(http.StatusInternalServerError, "Failed to update token to complete state")

		}

		return getToken(*token)
	}
}

func getToken(token tokenCode.TokenCode) (restApi.Response, error) {
	data := url.Values{}
	oauthUrl, err := url.ParseRequestURI(CognitoUrl)

	if err != nil {
		return restApi.BuildErrorResponse(http.StatusInternalServerError, "Error while parsing cognito url")
	}

	data.Set("redirect_uri", CallbackUri)
	data.Set("client_id", token.ClientId)
	data.Set("code", token.LoginCode)
	data.Set("grant_type", "authorization_code")
	data.Set("code_verifier", token.CodeVerifier)
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

	var tokenResponse CognitoTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return restApi.BuildErrorResponse(http.StatusInternalServerError, "Error decoding JSON response")
	}

	return restApi.ApiResponse(http.StatusCreated, APIResponse(tokenResponse))
}

func generateCodeVerifierAndChallenge() CodeChallenge {
	var (
		codeVerifier = uuid.New().String()
	)
	hasher := sha256.New()

	hasher.Write([]byte(codeVerifier))
	codeVerifierHash := hasher.Sum(nil)

	codeChallenge := base64.URLEncoding.EncodeToString(codeVerifierHash)
	codeChallenge = strings.TrimRight(codeChallenge, "=")

	return CodeChallenge{
		CodeChallenge: codeChallenge,
		CodeVerifier:  codeVerifier,
	}
}

// TODO: find a different way to validate it
func validateTokenCode(token tokenCode.TokenCode, clientId string) (*restApi.Response, error) {
	expiresAt, err := strconv.ParseInt(token.ExpireAt, 10, 64)

	if err != nil {
		response, err := restApi.BuildErrorResponse(http.StatusInternalServerError, "Impossible to convert expires at")

		return &response, err
	}

	var (
		isExpired       = time.Now().Unix() > expiresAt
		isStatusExpired = token.Status == "expired"
	)

	if isStatusExpired || isExpired {
		response, err := restApi.BuildErrorResponse(http.StatusGone, "Token is expired")

		if err == nil && !isStatusExpired {
			err = tokenCode.UpdateToken(token.DeviceCode, tokenCode.TokenCode{Status: "expired"})

			if err != nil {
				response, err = restApi.BuildErrorResponse(http.StatusInternalServerError, "Failed to update token")
			}
		}

		return &response, err
	}

	if token.ClientId != clientId {
		response, err := restApi.BuildErrorResponse(http.StatusForbidden, "The client ID does not match")

		return &response, err
	}

	lastCheckedAt, err := strconv.ParseInt(token.LastCheckedAt, 10, 64)

	if err != nil {
		response, err := restApi.BuildErrorResponse(http.StatusInternalServerError, "Impossible to convert last checked at")

		return &response, err
	}

	if time.Now().Unix() <= lastCheckedAt+PollingIntervalInSeconds {
		var response restApi.Response

		err = tokenCode.UpdateToken(token.DeviceCode, tokenCode.TokenCode{LastCheckedAt: strconv.FormatInt(time.Now().Unix(), 10)})

		if err != nil {
			response, err = restApi.BuildErrorResponse(http.StatusInternalServerError, "Failed to update token")

			return &response, err
		}

		response, err = restApi.BuildErrorResponse(http.StatusTooManyRequests, fmt.Sprintf("Respect the pooling interval of %d second(s)", PollingIntervalInSeconds))

		return &response, err
	}

	return nil, nil
}
