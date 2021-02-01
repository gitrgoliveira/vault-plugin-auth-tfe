package tfeauth

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *tfeAuthBackend) pathConfig() *framework.Path {
	return &framework.Path{
		Pattern: "config$",
		Fields: map[string]*framework.FieldSchema{
			"terraform_host": {
				Type:        framework.TypeString,
				Description: "TFE host (e.g. app.terraform.io)",
				Default:     "app.terraform.io",
			},
			"terraform_host_cert": {
				Type:        framework.TypeString,
				Description: "PEM encoded certificate used to validate the TFE host endpoint",
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "TFE CA Certificate",
				},
			},
			"organization": {
				Type:        framework.TypeString,
				Description: "TFE organization allowed to use this backend",
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigRead,
				Summary:  "Read the current authentication backend configuration.",
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigWrite,
				Summary:  "configure the auth backend",
			},
		},
	}
}

func (b *tfeAuthBackend) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	org := data.Get("organization").(string)
	if org == "" {
		return logical.ErrorResponse("no organization provided"), nil
	}
	host := data.Get("terraform_host").(string)
	if host == "" {
		return logical.ErrorResponse("no host provided"), nil
	}
	caCert := data.Get("terraform_host_cert").(string)

	config := &tfeConfig{
		Host:         host,
		CACert:       caCert,
		Organization: org,
	}

	entry, err := logical.StorageEntryJSON(configPath, config)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}
	return nil, nil
}

type tfeConfig struct {
	// Host is the url string for the TFE API
	Host string `json:"host"`
	// CACert is the CA Cert to use to call into the TFE API
	CACert string `json:"ca_cert"`
	// The organization autthorised to use this backend
	Organization string `json:"organization"`
}

func (b *tfeAuthBackend) pathConfigRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if config, err := b.config(ctx, req.Storage); err != nil {
		return nil, err
	} else if config == nil {
		return nil, nil
	} else {
		resp := &logical.Response{
			Data: map[string]interface{}{
				"terraform_host":      config.Host,
				"terraform_host_cert": config.CACert,
				"organization":        config.Organization,
			},
		}
		return resp, nil
	}
}
