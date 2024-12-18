package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	tokenCode "github.com/PBH-Tech/moonenv/lambdas/endpoints/auth"
	restApi "github.com/PBH-Tech/moonenv/lambdas/util/rest-api"
)

var characterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

type CodeChallenge struct {
	CodeChallenge string
	CodeVerifier  string
}

func randStringRunes(n int) string {
	b := make([]rune, n)

	for i := range b {
		b[i] = characterRunes[rand.Intn(len(characterRunes))]
	}

	return string(b)
}

func RequestSetOfToken(clientId string) (restApi.Response, error) {
	var (
		stateCodeLength  = 16
		deviceCodeLength = 16
		stateCode        = randStringRunes(stateCodeLength)
		deviceCode       = randStringRunes(deviceCodeLength)
		expiresIn        = 900 // 15 minutes
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
		Status:                  "authorization_pending",
		ExpireAt:                strconv.FormatInt(time.Now().Add(time.Duration(expiresIn)*time.Second).Unix(), 10),
		LastCheckedAt:           strconv.FormatInt(time.Now().Unix(), 10),
	})

	if err != nil {
		return restApi.ApiResponse(http.StatusBadRequest, err.Error())
	}

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

	if token.Status == "authorization_pending" || token.Status == "denied" {
		return restApi.ApiResponse(http.StatusBadRequest, token)
	} else {
		return restApi.ApiResponse(http.StatusOK, map[string]string{"message": "Authorized"})
	}
}

func generateCodeVerifierAndChallenge() CodeChallenge {
	var (
		codeVerifierLength = 32
		codeVerifier       = randStringRunes(codeVerifierLength)
	)
	hasher := sha256.New()

	hasher.Write([]byte(codeVerifier))
	challenge := hasher.Sum(nil)

	return CodeChallenge{
		CodeChallenge: base64.URLEncoding.EncodeToString(challenge),
		CodeVerifier:  codeVerifier,
	}
}

// TODO: find a different way to validate it
func validateTokenCode(token tokenCode.TokenCode, clientId string) (*restApi.Response, error) {
	expiresAt, err := strconv.ParseInt(token.ExpireAt, 10, 64)

	if err != nil {
		println("%s", err.Error())
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
