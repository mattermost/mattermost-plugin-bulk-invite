import { GlobalState } from '@mattermost/types/lib/store';
import {manifest} from './manifest';

import {ReducerState} from './reducers';

const getPluginState = (state: GlobalState): ReducerState => state[`plugins-${manifest.id}`] || {};

export const isBulkAddChannelModalVisible = (state: GlobalState) => getPluginState(state).bulkAddChannelModalVisible;

export const getBulkAddChannelModal = (state: GlobalState) => getPluginState(state).bulkAddChannelModal;
