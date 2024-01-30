package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/mattermost/mattermost-plugin-bulk-invite/server/api"
	"github.com/mattermost/mattermost-plugin-bulk-invite/server/engine"
	"github.com/mattermost/mattermost-plugin-bulk-invite/server/kvstore"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/pluginapi"

	root "github.com/mattermost/mattermost-plugin-bulk-invite"
)

var (
	Manifest model.Manifest = root.Manifest
)

const enableDebugRoute = false

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	// HTTP
	handler *api.Handler

	// botUserID the userID for the user of the bot, used to send messages to channels
	botUserID string

	// engine the engine to use on bulk operations
	engine *engine.Engine
}

func (p *Plugin) OnActivate() error {
	config := p.API.GetConfig()
	license := p.API.GetLicense()

	if !pluginapi.IsEnterpriseLicensedOrDevelopment(config, license) {
		return fmt.Errorf("this plugin requires an Enterprise license")
	}

	return nil
}

func (p *Plugin) OnDeactivate() error {
	return nil
}

// OnConfigurationChange is invoked when configuration changes may have been made.
func (p *Plugin) OnConfigurationChange() error {
	p.API.LogDebug("config change")
	var configuration = new(configuration)

	// Load the public configuration fields from the Mattermost server configuration.
	if err := p.API.LoadPluginConfiguration(configuration); err != nil {
		return fmt.Errorf("failed to load plugin configuration: %w", err)
	}

	p.setConfiguration(configuration)

	if err := p.ensureBot(); err != nil {
		return fmt.Errorf("error ensuring bot is present: %w", err)
	}

	lockStore := kvstore.NewLockStore(p.API)

	p.engine = engine.NewEngine(p.API, lockStore, p.botUserID)

	p.handler = api.NewHandler(p.API)
	api.Init(p.handler, p.engine, enableDebugRoute)

	return nil
}

func (p *Plugin) ServeHTTP(_ *plugin.Context, w http.ResponseWriter, req *http.Request) {
	p.handler.ServeHTTP(w, req)
}

// ensureBot ensures that the bot user is present in the system
func (p *Plugin) ensureBot() error {
	p.API.LogDebug("ensuring bot user is present")

	if p.botUserID != "" {
		return nil
	}

	botUser := &model.Bot{
		OwnerId:     Manifest.Id, // Workaround to support older server version affected by https://github.com/mattermost/mattermost-server/pull/21560
		Username:    "bulk-invite",
		DisplayName: "Bulk Invite",
		Description: "Bulk invite bot",
	}

	botUserID, err := p.API.EnsureBotUser(botUser)
	if err != nil {
		return fmt.Errorf("failed to ensure bot account: %w", err)
	}
	p.botUserID = botUserID
	return nil
}
