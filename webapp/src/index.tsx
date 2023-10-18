import {Store, Action} from 'redux';

import {GlobalState} from '@mattermost/types/lib/store';

import {PluginRegistry} from '@/types/mattermost-webapp';
import { bulkInvite, checkPost, openBulkInviteChannelModal } from './actions';
import { PluginId } from './plugin_id';
import { useEffect } from 'react';
import reducers from './reducers';
import BulkInviteChannelModal from './components/modals/bulk_invite_modal';

export default class Plugin {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars, @typescript-eslint/no-empty-function
    public async initialize(registry: PluginRegistry, store: Store<GlobalState, Action<Record<string, unknown>>>) {
        registry.registerReducer(reducers);

        const setup = async () => {
            registry.registerChannelHeaderMenuAction(
                "Bulk invite",
                async (channelID: string) => {
                    store.dispatch(openBulkInviteChannelModal(channelID));
                },
                () => { return true },
            );
            registry.registerRootComponent(BulkInviteChannelModal);
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
