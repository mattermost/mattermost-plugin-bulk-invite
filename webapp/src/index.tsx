import {Store, Action} from 'redux';

import {GlobalState} from '@mattermost/types/lib/store';

import {useEffect} from 'react';

import {PluginRegistry} from '@/types/mattermost-webapp';

import {openBulkAddChannelModal} from './actions';
import {manifest} from './manifest';
import reducers from './reducers';
import BulkAddChannelModal from './components/modals/bulk_add_channel_modal';

export default class Plugin {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars, @typescript-eslint/no-empty-function
    public async initialize(registry: PluginRegistry, store: Store<GlobalState, Action<Record<string, unknown>>>) {
        registry.registerReducer(reducers);

        const setup = async () => {
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
