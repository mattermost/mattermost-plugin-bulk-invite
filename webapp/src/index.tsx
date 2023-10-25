import {Store, Action} from 'redux';

import {GlobalState} from '@mattermost/types/lib/store';

import React, {useEffect} from 'react';

import {GlobalState as ReduxGlobalState} from 'mattermost-redux/types/store';

import {PluginRegistry} from '@/types/mattermost-webapp';

import {manifest} from './manifest';
import reducers, {PluginAction} from './reducers';
import BulkAddChannelModal from './components/modals/bulk_add_channel_modal';
import {setupClient} from './client';
import {openBulkAddChannelModal} from './actions';

export default class Plugin {
    public async initialize(registry: PluginRegistry, store: Store<GlobalState, PluginAction>) {
        registry.registerReducer(reducers);

        const setup = async () => {
            setupClient(store.getState() as any as ReduxGlobalState);

            registry.registerChannelHeaderMenuAction(
                'Bulk invite',
                async (channelID: string) => {
                    store.dispatch(openBulkAddChannelModal(channelID));
                },
            );

            registry.registerRootComponent(BulkAddChannelModal);
        };

        registry.registerRootComponent(() => <SetupUI setup={setup}/>);
    }
}

const SetupUI = ({setup}: { setup: () => Promise<void> }) => {
    useEffect(() => {
        setup();
    }, []);

    return null;
};

declare global {
    interface Window {
        registerPlugin(pluginId: string, plugin: Plugin): void
    }
}

window.registerPlugin(manifest.id, new Plugin());
