import {Store, Action} from 'redux';

import {GlobalState} from '@mattermost/types/lib/store';

import PluginRegistry from '@mattermost/webapp/plugins/registry';

import {useEffect} from 'react';

// import {Channel, ChannelType} from 'mattermost-redux/types/channels';

import {alwaysShow, openBulkAddChannelModal} from './actions';
import {PluginId} from './plugin_id';
import reducers from './reducers';
import BulkAddChannelModal from './components/modals/bulk_add_modal';

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

const SetupUI = ({setup}) => {
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

window.registerPlugin(PluginId, new Plugin());
