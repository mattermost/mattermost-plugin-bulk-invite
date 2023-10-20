import Client4 from 'mattermost-redux/client/client4';

import {getConfig} from 'mattermost-redux/selectors/entities/general';

import {Channel} from 'mattermost-redux/types/channels';

import {GlobalState} from 'mattermost-redux/types/store';

import {doFetchWithResponse} from './client';
import {getPluginServerRoute} from './selectors';
import {BulkAddChannelPayload} from './components/forms/bulk_add_channel_form';
import action_types from './action_types';

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

export type BulkAddChannelEventResponse = {data?: any; error?: string};

export const bulkAddToChannel = (payload: BulkAddChannelPayload) => async (dispatch, getState): Promise<BulkAddChannelEventResponse> => {
    const state = getState();
    const pluginServerRoute = getPluginServerRoute(state);

    const formData = new FormData();
    formData.append('channel_id', payload.channel_id);
    formData.append('add_to_team', String(payload.add_to_team).toLowerCase());
    formData.append('add_guests', String(payload.add_guests).toLowerCase());
    formData.append('file', payload.file);

    return doFetchWithResponse(`${pluginServerRoute}/handlers/channel_bulk_add`, {
        method: 'POST',
        body: formData,
    }).
        then((data) => {
            return {data};
        }).
        catch((response) => {
            const error = response.message?.error || 'An error occurred. Please check logs.';
            return {error};
        });
};

export const openBulkAddChannelModal = (channelId: string) => {
    return {
        type: action_types.OPEN_BULK_ADD_CHANNEL_MODAL,
        data: {
            channelId,
        },
    };
};

export const closeBulkAddChannelModal = () => {
    return {
        type: action_types.CLOSE_BULK_ADD_CHANNEL_MODAL,
    };
};

export type GetChannelResponse = {channel?: Channel | null; error?: string | null};

export const getChannelInfo = (channelId: string) => async (dispatch, getState): Promise<GetChannelResponse> => {
    const state = getState();
    const siteURL = getSiteURL(state);
    client.setUrl(siteURL);

    try {
        const channel = await client.getChannel(channelId);
        return {channel, error: null};
    } catch (e) {
        const error = e.message?.error || 'An error occurred while retrieving channel information.';
        return {channel: null, error};
    }
};
