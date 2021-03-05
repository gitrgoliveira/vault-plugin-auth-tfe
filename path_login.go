package tfeauth

import (
	"context"
	"errors"
	"fmt"

	log "github.com/hashicorp/go-hclog"
	tfe "github.com/hashicorp/go-tfe"
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

func (b *tfeAuthBackend) parseAndValidateLogin(role *roleStorageEntry, config *tfeAuthConfig,
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

func (t *tfeLogin) lookup(role *roleStorageEntry, config *tfeAuthConfig) error {

	clientConfig := &tfe.Config{
		Address: config.Host,
		Token:   t.AtlasToken,
	}

	ctx := context.Background()

	client, err := tfe.NewClient(clientConfig)
	if err != nil {
		msg := fmt.Sprintf("Error creating client for host %s with token %s -> %s", config.Host, t.AtlasToken, string(err.Error()))
		return fmt.Errorf(msg)
	}

	run, err := client.Runs.Read(ctx, t.RunID)
	if err != nil {
		msg := fmt.Sprintf("Error fetching RunID %s Info: %s", t.RunID, string(err.Error()))
		return fmt.Errorf(msg)
	}

	workspace, err := client.Workspaces.Read(ctx, config.Organization, t.Workspace)
	if err != nil {
		msg := fmt.Sprintf("Error fetching Workspace %s Info -> %s", t.Workspace, string(err.Error()))
		return fmt.Errorf(msg)
	}

	account, err := client.Users.ReadCurrent(ctx)
	if err != nil {
		msg := fmt.Sprintf("Error fetching Account Info for token %s -> %s", t.AtlasToken, string(err.Error()))
		return fmt.Errorf(msg)
	}

	msg := fmt.Sprintf("Run status is %s", run.Status)
	log.L().Info(msg, "info", nil)

	// Run must be active
	if run.Status != "applying" &&
		run.Status != "planning" {
		msg := fmt.Sprintf("Run ID %s status is %s. Expected planning or applying", run.ID, run.Status)
		return fmt.Errorf(msg)
	}

	// The Run must be related to the specified workspace
	if run.Workspace.ID != workspace.ID {
		msg := fmt.Sprintf("Workspace ID in Run (%s) and workspace ID (%s) mismatch", run.ID, workspace.ID)
		return fmt.Errorf(msg)
	}

	// The account must be a service account.
	if account.IsServiceAccount == false {
		msg := fmt.Sprintf("ATLAS Token must belong to a service account")
		return fmt.Errorf(msg)
	}

	log.L().Info(string(run.ID), "info", nil)
	log.L().Info(string(workspace.ID), "info", nil)
	log.L().Info(string(account.ID), "info", nil)

	return nil
}

const pathLoginHelpSyn = `Authenticates the current workspace run ID with Vault.`
const pathLoginHelpDesc = `
Authenticate current workspace run ID.
`
