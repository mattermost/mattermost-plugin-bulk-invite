import {Store, Action} from 'redux';

import {GlobalState} from '@mattermost/types/lib/store';

import React, {useEffect} from 'react';

import {GlobalState as ReduxGlobalState} from 'mattermost-redux/types/store';

import {getCurrentChannel} from 'mattermost-redux/selectors/entities/channels';
import Constants from 'mattermost-redux/constants/general';

import {PluginRegistry} from '@/types/mattermost-webapp';

import {manifest} from './manifest';
import reducers, {PluginAction} from './reducers';
import BulkAddChannelModal from './components/modals/bulk_add_channel_modal';
import {setupClient} from './client';
import {openBulkAddChannelModal} from './actions';

export default class Plugin {
    private setupUIFinished = false;

    public async initialize(registry: PluginRegistry, store: Store<GlobalState, PluginAction>) {
        const setup = async () => {
            setupClient(store.getState() as unknown as ReduxGlobalState);

            registry.registerChannelHeaderMenuAction(
                'Bulk invite',
                async (channelID: string) => {
                    store.dispatch(openBulkAddChannelModal(channelID));
                },
                () => {
                    const currentChannel = getCurrentChannel(store.getState() as unknown as ReduxGlobalState);
                    return ![Constants.DM_CHANNEL, Constants.GM_CHANNEL].includes(currentChannel.type);
                },
            );

            registry.registerRootComponent(BulkAddChannelModal);

            this.setupUIFinished = true;
        };

        registry.registerReducer(reducers);
        registry.registerRootComponent(() => (
            <SetupUI
                setup={setup}
                setupUIFinished={this.setupUIFinished}
            />
        ));
    }
}

const SetupUI = ({setup, setupUIFinished}: { setup: () => Promise<void>, setupUIFinished: boolean }) => {
    useEffect(() => {
        if (!setupUIFinished) {
            setup();
        }
    }, []);

    return null;
};

declare global {
    interface Window {
        registerPlugin(pluginId: string, plugin: Plugin): void
    }
}

window.registerPlugin(manifest.id, new Plugin());
