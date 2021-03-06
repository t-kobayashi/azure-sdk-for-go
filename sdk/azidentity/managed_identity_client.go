// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	imdsEndpoint = "http://169.254.169.254/metadata/identity/oauth2/token"
)

const (
	msiEndpointEnvironemntVariable = "MSI_ENDPOINT"
	msiSecretEnvironemntVariable   = "MSI_SECRET"
	appServiceMsiAPIVersion        = "2017-09-01"
	imdsAPIVersion                 = "2018-02-01"
)

type msiType int

const (
	msiTypeUnknown     msiType = 0
	msiTypeIMDS        msiType = 1
	msiTypeAppService  msiType = 2
	msiTypeCloudShell  msiType = 3
	msiTypeUnavailable msiType = 4
)

// managedIdentityClient provides the base for authenticating in managed identity environments
// This type includes an azcore.Pipeline and TokenCredentialOptions.
type managedIdentityClient struct {
	pipeline               azcore.Pipeline
	imdsAPIVersion         string
	imdsAvailableTimeoutMS time.Duration
	msiType                msiType
	endpoint               string
}

type wrappedNumber json.Number

func (n *wrappedNumber) UnmarshalJSON(b []byte) error {
	c := string(b)
	if c == "\"\"" {
		return nil
	}
	return json.Unmarshal(b, (*json.Number)(n))
}

var (
	defaultMSIOpts = newDefaultManagedIdentityOptions()
)

func newDefaultManagedIdentityOptions() *ManagedIdentityCredentialOptions {
	return &ManagedIdentityCredentialOptions{
		HTTPClient: azcore.DefaultHTTPClientTransport(),
	}
}

// newManagedIdentityClient creates a new instance of the ManagedIdentityClient with the ManagedIdentityCredentialOptions
// that are passed into it along with a default pipeline.
// options: ManagedIdentityCredentialOptions configure policies for the pipeline and the authority host that
// will be used to retrieve tokens and authenticate
func newManagedIdentityClient(options *ManagedIdentityCredentialOptions) *managedIdentityClient {
	logEnvVars()
	options = options.setDefaultValues()
	return &managedIdentityClient{
		pipeline:               newDefaultMSIPipeline(*options), // a pipeline that includes the specific requirements for MSI authentication, such as custom retry policy options
		imdsAPIVersion:         imdsAPIVersion,                  // this field will be set to whatever value exists in the constant and is used when creating requests to IMDS
		imdsAvailableTimeoutMS: 500,                             // we allow a timeout of 500 ms since the endpoint might be slow to respond
		msiType:                msiTypeUnknown,                  // when creating a new managedIdentityClient, the current MSI type is unknown and will be tested for and replaced once authenticate() is called from GetToken on the credential side
	}
}

// authenticate creates an authentication request for a Managed Identity and returns the resulting Access Token if successful.
// ctx: The current context for controlling the request lifetime.
// clientID: The client (application) ID of the service principal.
// scopes: The scopes required for the token.
func (c *managedIdentityClient) authenticate(ctx context.Context, clientID string, scopes []string) (*azcore.AccessToken, error) {
	currentMSI, err := c.getMSIType(ctx)
	if err != nil {
		return nil, err
	}
	// This condition should never be true since getMSIType returns an error in these cases
	// if MSI is unavailable or we were unable to determine the type return a nil access token
	if currentMSI == msiTypeUnavailable || currentMSI == msiTypeUnknown {
		return nil, &CredentialUnavailableError{CredentialType: "Managed Identity Credential", Message: "Please make sure you are running in a managed identity environment, such as a VM, Azure Functions, Cloud Shell, etc..."}
	}

	AT, err := c.sendAuthRequest(ctx, currentMSI, clientID, scopes)
	if err != nil {
		return nil, err
	}
	return AT, nil
}

func (c *managedIdentityClient) sendAuthRequest(ctx context.Context, msiType msiType, clientID string, scopes []string) (*azcore.AccessToken, error) {
	msg, err := c.createAuthRequest(ctx, msiType, clientID, scopes)
	if err != nil {
		return nil, err
	}

	resp, err := c.pipeline.Do(msg)
	if err != nil {
		return nil, err
	}

	if resp.HasStatusCode(successStatusCodes[:]...) {
		return c.createAccessToken(resp)
	}

	return nil, &AuthenticationFailedError{inner: newAADAuthenticationFailedError(resp)}
}

func (c *managedIdentityClient) createAccessToken(res *azcore.Response) (*azcore.AccessToken, error) {
	value := struct {
		// these are the only fields that we use
		Token        string        `json:"access_token,omitempty"`
		RefreshToken string        `json:"refresh_token,omitempty"`
		ExpiresIn    wrappedNumber `json:"expires_in,omitempty"` // this field should always return the number of seconds for which a token is valid
		ExpiresOn    string        `json:"expires_on,omitempty"` // the value returned in this field varies between a number and a date string
	}{}
	if err := res.UnmarshalAsJSON(&value); err != nil {
		return nil, fmt.Errorf("internal AccessToken: %w", err)
	}
	if value.ExpiresIn != "" {
		expiresIn, err := json.Number(value.ExpiresIn).Int64()
		if err != nil {
			return nil, err
		}
		return &azcore.AccessToken{Token: value.Token, ExpiresOn: time.Now().Add(time.Second * time.Duration(expiresIn)).UTC()}, nil
	}
	if expiresOn, err := strconv.Atoi(value.ExpiresOn); err == nil {
		return &azcore.AccessToken{Token: value.Token, ExpiresOn: time.Now().Add(time.Second * time.Duration(expiresOn)).UTC()}, nil
	}
	// this is the case when expires_on is a time string
	// this is the format of the string coming from the service
	if expiresOn, err := time.Parse("01/02/2006 15:04:05 PM +00:00", value.ExpiresOn); err == nil { // the date string specified in the layout param of time.Parse cannot be changed, Golang expects whatever layout to always signify January 2, 2006 at 3:04 PM
		eo := expiresOn.UTC()
		return &azcore.AccessToken{Token: value.Token, ExpiresOn: eo}, nil
	} else {
		return nil, err
	}
}

func (c *managedIdentityClient) createAuthRequest(ctx context.Context, msiType msiType, clientID string, scopes []string) (*azcore.Request, error) {
	switch msiType {
	case msiTypeIMDS:
		return c.createIMDSAuthRequest(ctx, scopes)
	case msiTypeAppService:
		return c.createAppServiceAuthRequest(ctx, clientID, scopes)
	case msiTypeCloudShell:
		return c.createCloudShellAuthRequest(ctx, clientID, scopes)
	default:
		errorMsg := ""
		switch msiType {
		case msiTypeUnavailable:
			errorMsg = "unavailable"
		default:
			errorMsg = "unknown"
		}

		return nil, &CredentialUnavailableError{CredentialType: "Managed Identity Credential", Message: "Make sure you are running in a valid Managed Identity Environment. Status: " + errorMsg}
	}
}

func (c *managedIdentityClient) createIMDSAuthRequest(ctx context.Context, scopes []string) (*azcore.Request, error) {
	request, err := azcore.NewRequest(ctx, http.MethodGet, c.endpoint)
	if err != nil {
		return nil, err
	}
	request.Header.Set(azcore.HeaderMetadata, "true")
	q := request.URL.Query()
	q.Add("api-version", c.imdsAPIVersion)
	q.Add("resource", strings.Join(scopes, " "))
	request.URL.RawQuery = q.Encode()
	return request, nil
}

func (c *managedIdentityClient) createAppServiceAuthRequest(ctx context.Context, clientID string, scopes []string) (*azcore.Request, error) {
	request, err := azcore.NewRequest(ctx, http.MethodGet, c.endpoint)
	if err != nil {
		return nil, err
	}
	request.Header.Set("secret", os.Getenv(msiSecretEnvironemntVariable))
	q := request.URL.Query()
	q.Add("api-version", appServiceMsiAPIVersion)
	q.Add("resource", strings.Join(scopes, " "))
	if clientID != "" {
		q.Add(qpClientID, clientID)
	}
	request.URL.RawQuery = q.Encode()
	return request, nil
}

func (c *managedIdentityClient) createCloudShellAuthRequest(ctx context.Context, clientID string, scopes []string) (*azcore.Request, error) {
	request, err := azcore.NewRequest(ctx, http.MethodPost, c.endpoint)
	if err != nil {
		return nil, err
	}
	request.Header.Set(azcore.HeaderMetadata, "true")
	data := url.Values{}
	data.Set("resource", strings.Join(scopes, " "))
	if clientID != "" {
		data.Set("client_id", clientID)
	}
	dataEncoded := data.Encode()
	body := azcore.NopCloser(strings.NewReader(dataEncoded))
	if err := request.SetBody(body, azcore.HeaderURLEncoded); err != nil {
		return nil, err
	}
	return request, nil
}

func (c *managedIdentityClient) getMSIType(ctx context.Context) (msiType, error) {
	if c.msiType == msiTypeUnknown { // if we haven't already determined the msi type
		if endpointEnvVar := os.Getenv(msiEndpointEnvironemntVariable); endpointEnvVar != "" { // if the env var MSI_ENDPOINT is set
			c.endpoint = endpointEnvVar
			if secretEnvVar := os.Getenv(msiSecretEnvironemntVariable); secretEnvVar != "" { // if BOTH the env vars MSI_ENDPOINT and MSI_SECRET are set the MsiType is AppService
				c.msiType = msiTypeAppService
			} else { // if ONLY the env var MSI_ENDPOINT is set the MsiType is CloudShell
				c.msiType = msiTypeCloudShell
			}
		} else if c.imdsAvailable(ctx) { // if MSI_ENDPOINT is NOT set AND the IMDS endpoint is available the MsiType is Imds. This will timeout after 500 milliseconds
			c.endpoint = imdsEndpoint
			c.msiType = msiTypeIMDS
		} else { // if MSI_ENDPOINT is NOT set and IMDS enpoint is not available ManagedIdentity is not available
			c.msiType = msiTypeUnavailable
			return c.msiType, &CredentialUnavailableError{CredentialType: "Managed Identity Credential", Message: "Make sure you are running in a valid Managed Identity Environment"}
		}
	}
	return c.msiType, nil
}

func (c *managedIdentityClient) imdsAvailable(ctx context.Context) bool {
	tempCtx, cancel := context.WithTimeout(ctx, c.imdsAvailableTimeoutMS*time.Millisecond)
	defer cancel()
	// this should never fail
	request, _ := azcore.NewRequest(tempCtx, http.MethodGet, imdsEndpoint)
	q := request.URL.Query()
	q.Add("api-version", c.imdsAPIVersion)
	request.URL.RawQuery = q.Encode()
	resp, err := c.pipeline.Do(request)
	if err == nil {
		resp.Drain()
	}
	return err == nil
}
