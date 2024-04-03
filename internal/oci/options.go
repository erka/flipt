package oci

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"go.flipt.io/flipt/internal/containers"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/registry/remote/auth"
)

var ErrNoAWSECRAuthorizationData = errors.New("no ecr authorization data provided")

// StoreOptions are used to configure call to NewStore
// This shouldn't be handled directory, instead use one of the function options
// e.g. WithBundleDir or WithCredentials
type StoreOptions struct {
	bundleDir       string
	manifestVersion oras.PackManifestVersion
	auth            credentialFunc
}

// WithCredentials configures username and password credentials used for authenticating
// with remote registries
func WithCredentials(kind, user, pass string) containers.Option[StoreOptions] {
	switch kind {
	case "aws-ecr":
		return WithAWSECRCredentials()
	default:
		return WithStaticCredentials(user, pass)
	}
}

// WithStaticCredentials configures username and password credentials used for authenticating
// with remote registries
func WithStaticCredentials(user, pass string) containers.Option[StoreOptions] {
	return func(so *StoreOptions) {
		so.auth = func(registry string) auth.CredentialFunc {
			return auth.StaticCredential(registry, auth.Credential{
				Username: user,
				Password: pass,
			})
		}
	}
}

// WithAWSECRCredentials configures username and password credentials used for authenticating
// with remote registries
func WithAWSECRCredentials() containers.Option[StoreOptions] {
	return func(so *StoreOptions) {
		so.auth = func(registry string) auth.CredentialFunc {
			// AWS default private registry is {aws_account_id}.dkr.ecr.{region}.amazonaws.com
			return awsECRCredential
		}
	}
}

// WithManifestVersion configures what OCI Manifest version to build the bundle.
func WithManifestVersion(version oras.PackManifestVersion) containers.Option[StoreOptions] {
	return func(s *StoreOptions) {
		s.manifestVersion = version
	}
}

func awsECRCredential(ctx context.Context, hostport string) (auth.Credential, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return auth.EmptyCredential, err
	}
	client := ecr.NewFromConfig(cfg)
	response, err := client.GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return auth.EmptyCredential, err
	}
	if len(response.AuthorizationData) == 0 {
		return auth.EmptyCredential, ErrNoAWSECRAuthorizationData
	}
	return credentialsFromECRAuthorizationToken(response.AuthorizationData[0].AuthorizationToken)
}

// credentialsFromECRAuthorizationToken converts the raw AWS ECR token to oras credential
func credentialsFromECRAuthorizationToken(token *string) (auth.Credential, error) {
	if token == nil {
		return auth.EmptyCredential, auth.ErrBasicCredentialNotFound
	}

	output, err := base64.StdEncoding.DecodeString(*token)
	if err != nil {
		return auth.EmptyCredential, err
	}

	userpass := strings.SplitN(string(output), ":", 2)
	if len(userpass) != 2 {
		return auth.EmptyCredential, auth.ErrBasicCredentialNotFound
	}

	return auth.Credential{
		Username: userpass[0],
		Password: userpass[1],
	}, nil
}
