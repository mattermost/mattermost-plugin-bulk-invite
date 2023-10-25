import {Reducer, combineReducers} from 'redux';

import action_types from './action_types';

export type PluginAction = {
    type: string;
    data?: object;
};

type BulkAddChannelModalActionData = {
    channelId?: string;
};

export type BulkAddChannelModalAction = {
    type: string;
    data?: BulkAddChannelModalActionData;
};

const bulkAddChannelModalVisible: Reducer<boolean, BulkAddChannelModalAction> = (state = false, action) => {
    switch (action.type) {
    case action_types.OPEN_BULK_ADD_CHANNEL_MODAL:
        return true;
    case action_types.CLOSE_BULK_ADD_CHANNEL_MODAL:
        return false;
    default:
        return state;
    }
};

const bulkAddChannelModal: Reducer<BulkAddChannelModalActionData, BulkAddChannelModalAction> = (state = {}, action) => {
    switch (action.type) {
    case action_types.OPEN_BULK_ADD_CHANNEL_MODAL:
        return {
            channelId: action.data?.channelId,
        };
    case action_types.CLOSE_BULK_ADD_CHANNEL_MODAL:
        return {};
    default:
        return state;
    }
};

export default combineReducers({
    bulkAddChannelModal,
    bulkAddChannelModalVisible,
});

export type ReducerState = {
    bulkAddChannelModalVisible: boolean;
    bulkAddChannelModal: {
        channelId: string;
    } | null;
}
