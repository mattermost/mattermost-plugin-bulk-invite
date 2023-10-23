import {manifest} from './manifest';

import {ReducerState} from './reducers';

const getPluginState = (state): ReducerState => state['plugins-' + manifest.id] || {};

export const isBulkAddChannelModalVisible = (state) => getPluginState(state).bulkAddChannelModalVisible;

export const getBulkAddChannelModal = (state) => getPluginState(state).bulkAddChannelModal;
