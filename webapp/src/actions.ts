import Client4 from 'mattermost-redux/client/client4';

import { doFetchWithResponse } from "./client";
import { getPluginServerRoute } from "./selectors";
import { getConfig } from 'mattermost-redux/selectors/entities/general';
import { BulkInvitePayload } from './components/forms/bulk_invite_channel_form';
import action_types from './action_types';
import { Channel } from 'mattermost-redux/types/channels';

const client = new Client4();

export const getSiteURL = (state: GlobalState): string => {
    const config = getConfig(state);

    let basePath = '';
    if (config && config.SiteURL) {
        basePath = new URL(config.SiteURL).pathname;

        if (basePath && basePath[basePath.length - 1] === '/') {
            basePath = basePath.substring(0, basePath.length - 1);
        }
    }

    return basePath;
};

export const alwaysShow = (postId: string): boolean => {
    return true
}

export type BulkInviteChannelEventResponse = {data?: any; error?: string};

export const bulkInviteToChannel = (payload: BulkInvitePayload) => async (dispatch, getState): Promise<BulkInviteChannelEventResponse> => {
    const state = getState();
    const pluginServerRoute = getPluginServerRoute(state);

    return doFetchWithResponse(`${pluginServerRoute}/api/v1/bulk_invite`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
    }).
        then((data) => {
            return {data};
        }).
        catch((response) => {
            const error = response.message?.error || 'An error occurred. Please check logs.';
            return {error};
        });
};

export const openBulkInviteChannelModal = (channelId: string) => {
    return {
        type: action_types.OPEN_BULK_INVITE_CHANNEL_MODAL,
        data: {
            channelId,
        },
    };
};

export const closeBulkInviteChannelModal = () => {
    return {
        type: action_types.CLOSE_BULK_INVITE_CHANNEL_MODAL,
    };
};

export type GetChannelResponse = {channel?: Channel | null; error?: string | null};

export const getChannelInfo = (channelId: string) => async (dispatch, getState): Promise<GetChannelResponse> => {
    const state = getState();
    const siteURL = getSiteURL(state);
    client.setUrl(siteURL);

    try {
        const channel = await client.getChannel(channelId);
        return {channel: channel, error: null};
    } catch (e) {
        const error = e.message?.error || 'An error occurred while retrieving channel information.';
        return {channel: null, error};
    }
};
