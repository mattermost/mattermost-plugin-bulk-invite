package engine

import (
	"context"
	"sync"
	"testing"

	"github.com/mattermost/mattermost-plugin-bulk-invite/server/kvstore"
	"github.com/mattermost/mattermost-plugin-bulk-invite/server/mocks"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newValidEmptyConfig() *Config {
	return &Config{
		ChannelID: "test",
		UserID:    "user-id",
		Users:     []AddUser{},
		AddToTeam: false,
	}
}

type engineTestHelper struct {
	ctrl *gomock.Controller

	API *plugintest.API
	KV  kvstore.LockStore
}

func (h *engineTestHelper) finish() {
	h.ctrl.Finish()
}

func newEngineTestHelper(t *testing.T) *engineTestHelper {
	ctrl := gomock.NewController(t)
	return &engineTestHelper{
		ctrl: ctrl,
		API:  plugintest.NewAPI(t),
		KV:   mocks.NewMockLockStore(ctrl),
	}
}

func TestEngineStartJobErrors(t *testing.T) {
	t.Run("Lock", func(t *testing.T) {
		t.Run("Locked channel should fail", func(t *testing.T) {
			th := newEngineTestHelper(t)
			defer th.finish()
			engine := NewEngine(th.API, th.KV, "bot-user-id")

			cfg := newValidEmptyConfig()
			th.KV.(*mocks.MockLockStore).EXPECT().IsLocked(cfg.ChannelID).Return(true)

			err := engine.StartJob(context.TODO(), cfg)
			require.Error(t, err)
		})
	})

	t.Run("GetChannel errors", func(t *testing.T) {
		th := newEngineTestHelper(t)
		defer th.finish()
		engine := NewEngine(th.API, th.KV, "bot-user-id")

		cfg := newValidEmptyConfig()
		th.KV.(*mocks.MockLockStore).EXPECT().IsLocked(cfg.ChannelID).Return(false)
		appErr := model.AppError{
			Where:         "",
			DetailedError: "some error",
		}
		th.API.On("LogError", "error getting channnel information", "channel_id", cfg.ChannelID, "err", appErr.Error())
		th.API.On("GetChannel", cfg.ChannelID).Return(nil, &appErr)

		err := engine.StartJob(context.TODO(), cfg)
		require.Error(t, err)
	})

	t.Run("Channel Type", func(t *testing.T) {
		t.Run("Group should fail", func(t *testing.T) {
			th := newEngineTestHelper(t)
			defer th.finish()
			engine := NewEngine(th.API, th.KV, "bot-user-id")

			cfg := newValidEmptyConfig()
			th.KV.(*mocks.MockLockStore).EXPECT().IsLocked(cfg.ChannelID).Return(false)
			th.API.On("GetChannel", cfg.ChannelID).Return(&model.Channel{
				Type: model.ChannelTypeGroup,
			}, nil)

			err := engine.StartJob(context.TODO(), cfg)
			require.Error(t, err)
		})

		t.Run("DM should fail", func(t *testing.T) {
			th := newEngineTestHelper(t)
			defer th.finish()
			engine := NewEngine(th.API, th.KV, "bot-user-id")

			cfg := newValidEmptyConfig()
			th.KV.(*mocks.MockLockStore).EXPECT().IsLocked(cfg.ChannelID).Return(false)
			th.API.On("GetChannel", cfg.ChannelID).Return(&model.Channel{
				Type: model.ChannelTypeDirect,
			}, nil)

			err := engine.StartJob(context.TODO(), cfg)
			require.Error(t, err)
		})
	})

	t.Run("User Permissions", func(t *testing.T) {
		t.Run("private channel without permissions should fail", func(t *testing.T) {
			th := newEngineTestHelper(t)
			defer th.finish()
			engine := NewEngine(th.API, th.KV, "bot-user-id")

			cfg := newValidEmptyConfig()
			th.KV.(*mocks.MockLockStore).EXPECT().IsLocked(cfg.ChannelID).Return(false)

			th.API.On("GetChannel", cfg.ChannelID).Return(&model.Channel{
				Type: model.ChannelTypePrivate,
			}, nil)
			th.API.On("HasPermissionToChannel", cfg.UserID, cfg.ChannelID, model.PermissionManagePrivateChannelMembers).Return(false)

			err := engine.StartJob(context.TODO(), cfg)
			require.Error(t, err)
		})

		t.Run("public channel without permissions should fail", func(t *testing.T) {
			th := newEngineTestHelper(t)
			defer th.finish()
			engine := NewEngine(th.API, th.KV, "bot-user-id")

			cfg := newValidEmptyConfig()
			th.KV.(*mocks.MockLockStore).EXPECT().IsLocked(cfg.ChannelID).Return(false)

			th.API.On("GetChannel", cfg.ChannelID).Return(&model.Channel{
				Type: model.ChannelTypeOpen,
			}, nil)
			th.API.On("HasPermissionToChannel", cfg.UserID, cfg.ChannelID, model.PermissionManagePublicChannelMembers).Return(false)

			err := engine.StartJob(context.TODO(), cfg)
			require.Error(t, err)
		})

		t.Run("Add to team without permissions should fail", func(t *testing.T) {
			th := newEngineTestHelper(t)
			defer th.finish()
			engine := NewEngine(th.API, th.KV, "bot-user-id")

			cfg := newValidEmptyConfig()
			cfg.AddToTeam = true
			th.KV.(*mocks.MockLockStore).EXPECT().IsLocked(cfg.ChannelID).Return(false)

			th.API.On("GetChannel", cfg.ChannelID).Return(&model.Channel{
				Type:   model.ChannelTypeOpen,
				TeamId: "team-id",
			}, nil)
			th.API.On("HasPermissionToChannel", cfg.UserID, cfg.ChannelID, model.PermissionManagePublicChannelMembers).Return(true)
			th.API.On("HasPermissionToTeam", cfg.UserID, "team-id", model.PermissionAddUserToTeam).Return(false)

			err := engine.StartJob(context.TODO(), cfg)
			require.Error(t, err)
		})
	})
}

func TestStartJobSuccess(t *testing.T) {
	th := newEngineTestHelper(t)
	defer th.finish()
	engine := NewEngine(th.API, th.KV, "bot-user-id")

	cfg := newValidEmptyConfig()

	th.KV.(*mocks.MockLockStore).EXPECT().IsLocked(cfg.ChannelID).Return(false)
	th.API.On("GetChannel", cfg.ChannelID).Return(&model.Channel{
		Type:   model.ChannelTypeOpen,
		TeamId: "team-id",
	}, nil)
	th.API.On("HasPermissionToChannel", cfg.UserID, cfg.ChannelID, model.PermissionManagePublicChannelMembers).Return(true)
	th.KV.(*mocks.MockLockStore).EXPECT().Lock(cfg.ChannelID).Return(nil)

	th.API.On("GetUser", cfg.UserID).Return(&model.User{
		Id:       cfg.UserID,
		Username: "username",
	}, nil)
	th.API.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(&model.Post{}, nil)
	th.KV.(*mocks.MockLockStore).EXPECT().Unlock(cfg.ChannelID).Return(nil)

	wg := sync.WaitGroup{}
	wg.Add(1)

	engine.SetOnFinish(func() {
		wg.Done()
	})
	err := engine.StartJob(context.Background(), cfg)
	require.Nil(t, err)

	// Wait for goroutine to finish
	wg.Wait()
}
