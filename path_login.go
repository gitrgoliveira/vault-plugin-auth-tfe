package tfeauth

import (
	"context"
	"errors"
	"fmt"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/cidrutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathLogin(b *tfeAuthBackend) *framework.Path {
	return &framework.Path{
		Pattern: "login$",
		Fields: map[string]*framework.FieldSchema{
			"role": {
				Type:        framework.TypeString,
				Description: `Name of the role against which the login is being attempted. This field is required`,
			},
			"workspace": {
				Type:        framework.TypeString,
				Description: "Name of the workspace that is loging in",
			},
			"run-id": {
				Type:        framework.TypeString,
				Description: "TFC_RUN_ID or ATLAS_RUN_ID of the current active run",
			},
			"atlas-token": {
				Type:        framework.TypeString,
				Description: "The ATLAS_TOKEN environment variable",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation:         b.pathLogin,
			logical.AliasLookaheadOperation: b.aliasLookahead,
		},

		HelpSynopsis:    pathLoginHelpSyn,
		HelpDescription: pathLoginHelpDesc,
	}
}

func (b *tfeAuthBackend) pathLogin(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role").(string)
	if len(roleName) == 0 {
		return logical.ErrorResponse("missing role"), nil
	}

	workspaceStr := data.Get("workspace").(string)
	if len(workspaceStr) == 0 {
		return logical.ErrorResponse("missing workspace"), nil
	}

	runIDStr := data.Get("run-id").(string)
	if len(runIDStr) == 0 {
		return logical.ErrorResponse("missing run-id"), nil
	}

	atlasTokenStr := data.Get("atlas-token").(string)
	if len(atlasTokenStr) == 0 {
		return logical.ErrorResponse("missing atlas-token"), nil
	}

	b.l.RLock()
	defer b.l.RUnlock()

	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("invalid role name \"%s\"", roleName)), nil
	}

	// Check for a CIDR match.
	if len(role.TokenBoundCIDRs) > 0 {
		if req.Connection == nil {
			b.Logger().Warn("token bound CIDRs found but no connection information available for validation")
			return nil, logical.ErrPermissionDenied
		}
		if !cidrutil.RemoteAddrIsOk(req.Connection.RemoteAddr, role.TokenBoundCIDRs) {
			return nil, logical.ErrPermissionDenied
		}
	}

	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, errors.New("could not load backend configuration")
	}

	tfeLogin, err := b.parseAndValidateLogin(role, config,
		workspaceStr, runIDStr, atlasTokenStr)
	if err != nil {
		return nil, err
	}

	err = tfeLogin.lookup(role, config)
	if err != nil {
		b.Logger().Error(`login unauthorized due to: ` + err.Error())
		return nil, logical.ErrPermissionDenied
	}

	auth := &logical.Auth{
		Alias: &logical.Alias{
			Name: fmt.Sprintf("%s/%s", config.Organization, tfeLogin.Workspace),
			Metadata: map[string]string{
				"Workspace":    tfeLogin.Workspace,
				"Organization": config.Organization,
			},
		},
		InternalData: map[string]interface{}{
			"role": roleName,
		},
		Metadata: map[string]string{
			"Workspace":    tfeLogin.Workspace,
			"Organization": config.Organization,
			"role":         roleName,
		},
		DisplayName: fmt.Sprintf("%s/%s", config.Organization, tfeLogin.Workspace),
	}

	role.PopulateTokenAuth(auth)

	return &logical.Response{
		Auth: auth,
	}, nil
}

func (b *tfeAuthBackend) aliasLookahead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	workspaceStr := data.Get("workspace").(string)
	if len(workspaceStr) == 0 {
		return logical.ErrorResponse("missing workspace"), nil
	}

	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, errors.New("could not load backend configuration")
	}

	return &logical.Response{
		Auth: &logical.Auth{
			Alias: &logical.Alias{
				Name: fmt.Sprintf("%s/%s", config.Organization, workspaceStr),
			},
		},
	}, nil
}

func (b *tfeAuthBackend) parseAndValidateLogin(role *roleStorageEntry, config *tfeConfig,
	workspace string, runID string, atlasToken string) (*tfeLogin, error) {

	if len(role.Workspaces) > 1 || role.Workspaces[0] != "*" {
		if !strutil.StrListContainsGlob(role.Workspaces, workspace) {
			return nil, errors.New("workspace not authorized")
		}
	}

	login := &tfeLogin{}
	login.Workspace = workspace
	login.RunID = runID
	login.AtlasToken = atlasToken

	return login, nil
}

type tfeLogin struct {
	Workspace  string `mapstructure:"workspace"`
	RunID      string `mapstructure:"run-id"`
	AtlasToken string `mapstructure:"atlas-token"`
}

func (t *tfeLogin) lookup(role *roleStorageEntry, config *tfeConfig) error {

	run, err := fetchRunInfo(t, config)
	if err != nil {
		msg := fmt.Sprintf("Error fetching Run Info: %s", string(err.Error()))
		return fmt.Errorf(msg)
	}

	workspace, err := fetchWorkspaceInfo(t, config)
	if err != nil {
		msg := fmt.Sprintf("Error fetching Workspace Info: %s", string(err.Error()))
		return fmt.Errorf(msg)
	}

	account, err := fetchAccountInfo(t, config)
	if err != nil {
		msg := fmt.Sprintf("Error fetching Account Info: %s", string(err.Error()))
		return fmt.Errorf(msg)
	}

	msg := fmt.Sprintf("Run status is %s", run.Data.Attributes.Status)
	log.L().Info(msg, "info", nil)

	// Run must be active
	if run.Data.Attributes.Status != "applying" &&
		run.Data.Attributes.Status != "planning" {
		msg := fmt.Sprintf("Run ID status is %s. Expected planning or applying", run.Data.Attributes.Status)
		return fmt.Errorf(msg)
	}

	if run.Data.Relationships.Workspace.Data.ID != workspace.Data.ID {
		msg := fmt.Sprintf("Workspace ID in Run (%s) and workspace ID (%s) mismatch", run.Data.Relationships.Workspace.Data.ID, workspace.Data.ID)
		return fmt.Errorf(msg)
	}

	log.L().Info(string(run.Data.ID), "info", nil)
	log.L().Info(string(workspace.Data.ID), "info", nil)
	log.L().Info(string(account.Data.ID), "info", nil)

	return nil
}

const pathLoginHelpSyn = `Authenticates the current workspace run ID with Vault.`
const pathLoginHelpDesc = `
Authenticate current workspace run ID.
`
