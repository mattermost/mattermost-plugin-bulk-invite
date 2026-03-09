import {Client4} from 'mattermost-redux/client';
import {Options} from '@mattermost/types/client4';

import {GlobalState} from '@mattermost/types/store';

class ClientError extends Error {
    url: string;
    status_code: number;

    constructor(baseUrl: string, data: {message: string; status_code: number; url: string}) {
        super(data.message);
        this.url = data.url;
        this.status_code = data.status_code;
    }
}

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
