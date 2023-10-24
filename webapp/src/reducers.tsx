import {AnyAction, combineReducers} from 'redux';

import action_types from './action_types';

export const openBulkInviteChannelModal = (channelId: string) => {
    return {
        type: action_types.OPEN_BULK_ADD_CHANNEL_MODAL,
        data: {
            channelId,
        },
    };
};

const bulkAddChannelModalVisible = (state = false, action: AnyAction) => {
    switch (action.type) {
    case action_types.OPEN_BULK_ADD_CHANNEL_MODAL:
        return true;
    case action_types.CLOSE_BULK_ADD_CHANNEL_MODAL:
        return false;
    default:
        return state;
    }
};

const bulkAddChannelModal = (state = false, action: AnyAction) => {
    switch (action.type) {
    case action_types.OPEN_BULK_ADD_CHANNEL_MODAL:
        return {
            channelId: action.data.channelId,
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
