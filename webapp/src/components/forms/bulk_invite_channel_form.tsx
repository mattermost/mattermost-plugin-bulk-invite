import React, {useCallback, useState} from 'react';
import {useSelector, useDispatch} from 'react-redux';

import {Modal} from 'react-bootstrap';

import {getTheme} from 'mattermost-redux/selectors/entities/preferences';

import FormButton from '../form_button';
import Loading from '../loading';
import FormComponent from '@/components/form_component';
import {BulkInviteChannelEventResponse, GetChannelResponse, bulkInviteToChannel, getChannelInfo} from '@/actions';

import './bulk_invite_channel_form.scss';
import {Props as FormComponentV2Props, FormComponentV2} from '../form_component_v2';
import { getBulkInviteChannelModal } from '@/selectors';
import { Channel } from 'mattermost-redux/types/channels';

type Props = {
    close: (e?: Event) => void;
};

export type BulkInvitePayload = {
    invite_to_team: boolean;
    file?: string
    users: string[];
    channel_id: string;
}

export default function BulkInviteChannelForm(props: Props) {
    const [storedError, setStoredError] = useState('');
    const [submitting, setSubmitting] = useState(false);
    const [loading, setLoading] = useState(false);
    const [channelName, setChannelName] = useState('');

    const dispatch = useDispatch();

    const modalProps = useSelector(getBulkInviteChannelModal);
    if (modalProps === null || modalProps.channelId === null) {
        return null;
    }

    const [formValues, setFormValues] = useState<BulkInvitePayload>({
        invite_to_team: false,
        users: [],
        channel_id: modalProps.channelId,
    });

    const loadChannelInfo = useCallback(async (channelId: string): Promise<Channel> => {
        const response = (await dispatch(getChannelInfo(channelId)) as unknown as GetChannelResponse);

        if (response.error) {
            setStoredError(response.error);
            return [];
        }

        setStoredError('');

        return response.channel;
    }, []);

    if (!channelName) {
        loadChannelInfo(modalProps.channelId).then((channel) => {
            setChannelName(channel.display_name)
        });
    }


    const setFormValue = <Key extends keyof BulkInvitePayload>(name: Key, value: BulkInvitePayload[Key]) => {
        setFormValues((values: BulkInvitePayload) => ({
            ...values,
            [name]: value,
        }));
    };

    const handleClose = (e?: Event) => {
        if (e && e.preventDefault) {
            e.preventDefault();
        }

        props.close();
    };

    const handleError = (error: string) => {
        setStoredError(error);
        setSubmitting(false);
    };

    const handleSubmit = async (e?: React.FormEvent) => {
        if (e && e.preventDefault) {
            e.preventDefault();
        }

        setSubmitting(true);

        const response = (await dispatch(bulkInviteToChannel(formValues))) as BulkInviteChannelEventResponse;
        if (response.error) {
            handleError(response.error);
            return;
        }

        handleClose();
    };

    const disableSubmit = false;
    const footer = (
        <React.Fragment>
            <FormButton
                type='button'
                btnClass='btn btn-tertiary'
                defaultMessage='Cancel'
                onClick={handleClose}
            />
            <FormButton
                id='submit-button'
                type='submit'
                btnClass='btn btn-primary'
                saving={submitting}
                disabled={disableSubmit}
            >
                {'Invite'}
            </FormButton>
        </React.Fragment>
    );

    let form;
    if (loading) {
        form = <Loading/>;
    } else {
        form = (
            <ActualForm
                formValues={formValues}
                setFormValue={setFormValue}
            />
        );
    }

    let error;
    if (storedError) {
        error = (
            <p className='alert alert-danger'>
                <i
                    className='fa fa-warning'
                    title='Warning Icon'
                />
                <span>{storedError}</span>
            </p>
        );
    }

    return (
        <form
            role='form'
            onSubmit={handleSubmit}
        >
            <Modal.Body >
                <div className='channel-invite__header'>
                    <h1>{'Bulk invite to ' + channelName}</h1>
                </div>
                {error}
                {form}
            </Modal.Body>
            <Modal.Footer>
                {footer}
            </Modal.Footer>
        </form>
    );
}

type ActualFormProps = {
    formValues: BulkInvitePayload;
    setFormValue: <Key extends keyof BulkInvitePayload>(name: Key, value: BulkInvitePayload[Key]) => Promise<{ error?: string }>;
}

const ActualForm = (props: ActualFormProps) => {
    const {formValues, setFormValue} = props;

    const theme = useSelector(getTheme);

    const components: FormComponentV2Props[] = [
        {
            label: 'Bulk invite file (.JSON format)',
            required: true,
            helpText: <div><a>Download a template</a> to ensure your file formatting is correct.</div>,
            element: (
                <input
                    id='bulk-invite-channel-file'
                    onChange={(e) => setFormValue('file', e.target.value)}
                    type='file'
                />
            ),
        },
        {
            label: 'Invite members to the team if they don’t belong to it',
            required: false,
            element: (
                <input
                    id='bulk-invite-channel-invite-to-team'
                    onChange={(e) => {
                        setFormValue('invite_to_team', e.target.checked)
                    }}
                    value={String(formValues.invite_to_team)}
                    type='checkbox'
                />
            ),
        },
    ];

    return (
        <div className='bulk-invite-channel-form'>
            {components.map((c) => (
                <FormComponentV2
                    {...c}
                    key={c.element.props.id}
                ></FormComponentV2>
            ))}
        </div>
    );
};
