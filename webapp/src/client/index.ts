import {Client4} from 'mattermost-redux/client';
import {Options} from 'mattermost-redux/types/client4';
import {ClientError} from 'mattermost-redux/client/client4';

import {getConfig} from 'mattermost-redux/selectors/entities/general';
import {GlobalState} from 'mattermost-redux/types/store';

import {getSiteURL} from '@/actions';

export const client = Client4;

export const doFetch = async (url: string, options: Options) => {
    const {data} = await doFetchWithResponse(url, options);

    return data;
};

export const doFetchWithResponse = async (url: string, options: Options = {}) => {
    const response = await fetch(url, Client4.getOptions(options));

    const data = await response.json();

    if (response.ok) {
        return {
            response,
            data,
        };
    }

    throw new ClientError(Client4.url, {
        message: data || '',
        status_code: response.status,
        url,
    });
};

export const doFormFetchWithResponse = async (url: string, options: Options = {}) => {
    const response = await fetch(url, Client4.getOptions(options));

    const data = await response.json();

    if (response.ok) {
        return {
            response,
            data,
        };
    }

    throw new ClientError(Client4.url, {
        message: data || '',
        status_code: response.status,
        url,
    });
};

export const setupClient = (state: GlobalState) => {
    client.setUrl(getSiteURL(state));
};
