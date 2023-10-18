import { combineReducers } from "redux";
import action_types from "./action_types";

export const openBulkInviteChannelModal = (channelId: string) => {
    return {
        type: action_types.OPEN_BULK_INVITE_CHANNEL_MODAL,
        data: {
            channelId,
        },
    };
};

const bulkInviteChannelModalVisible = (state = false, action) => {
    switch (action.type) {
    case action_types.OPEN_BULK_INVITE_CHANNEL_MODAL:
        return true;
    case action_types.CLOSE_BULK_INVITE_CHANNEL_MODAL:
        return false;
    default:
        return state;
    }
};


const bulkInviteChannelModal = (state = false, action) => {
    switch (action.type) {
    case action_types.OPEN_BULK_INVITE_CHANNEL_MODAL:
        return {
            channelId: action.data.channelId,
        };
    case action_types.CLOSE_BULK_INVITE_CHANNEL_MODAL:
        return {};
    default:
        return state;
    }
};

export default combineReducers({
    bulkInviteChannelModal,
    bulkInviteChannelModalVisible,
});

export type ReducerState = {
    bulkInviteChannelModalVisible: boolean;
    bulkInviteChannelModal: {
        channelId: string;
    } | null;
}
