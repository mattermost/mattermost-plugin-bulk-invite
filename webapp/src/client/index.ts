import {Client4} from 'mattermost-redux/client';
import {Options} from '@mattermost/types/client4';

import {GlobalState} from '@mattermost/types/store';

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

    throw new Error(data || '');
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

    throw new Error(data || '');
};

export const setupClient = (state: GlobalState) => {
    client.setUrl(getSiteURL(state));
};
