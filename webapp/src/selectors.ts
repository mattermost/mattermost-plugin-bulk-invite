import {GlobalState} from '@mattermost/types/lib/store';

import {manifest} from './manifest';

import {ReducerState} from './reducers';

const pluginStateProperty = `plugins-${manifest.id}`;
type GlobalStateWithPlugin = GlobalState & {[prop in typeof pluginStateProperty]: ReducerState};

const getPluginState = (state: GlobalStateWithPlugin): ReducerState => state[pluginStateProperty] || {};

export const isBulkAddChannelModalVisible = (state: GlobalStateWithPlugin) => getPluginState(state).bulkAddChannelModalVisible;

export const getBulkAddChannelModal = (state: GlobalStateWithPlugin) => getPluginState(state).bulkAddChannelModal;
