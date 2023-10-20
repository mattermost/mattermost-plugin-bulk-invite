import React, {useCallback, useState} from 'react';
import {useSelector, useDispatch} from 'react-redux';

import {Modal} from 'react-bootstrap';

import {Channel} from 'mattermost-redux/types/channels';

import FormButton from '../form_button';
import {BulkAddChannelEventResponse, GetChannelResponse, bulkAddToChannel, getChannelInfo} from '@/actions';

import {getBulkAddChannelModal} from '@/selectors';
import {Props as FormComponentProps, FormComponent} from '../form_component';

import './bulk_add_channel_form.scss';

type Props = {
    close: (e?: Event) => void;
};

export type BulkAddChannelPayload = {
    add_to_team: boolean;
    add_guests: boolean;
    file?: File
    users: string[];
    channel_id: string;
}

export default function BulkAddChannelForm(props: Props) {
    const [storedError, setStoredError] = useState('');
    const [submitting, setSubmitting] = useState(false);
    const [channelName, setChannelName] = useState('');

    const dispatch = useDispatch();

    const modalProps = useSelector(getBulkAddChannelModal);
    if (modalProps === null || modalProps.channelId === null) {
        return null;
    }

    const [formValues, setFormValues] = useState<BulkAddChannelPayload>({
        add_to_team: false,
        add_guests: false,
        users: [],
        channel_id: modalProps.channelId,
    });

    const loadChannelInfo = useCallback(async (channelId: string): Promise<Channel | null> => {
        const response = (await dispatch(getChannelInfo(channelId)) as unknown as GetChannelResponse);

        if (response.error) {
            setStoredError(response.error);
            return null;
        }

        setStoredError('');

        return response.channel;
    }, []);

    if (!channelName) {
        loadChannelInfo(modalProps.channelId).then((channel) => {
            if (!channel) setChannelName(channel.display_name);
        });
    }

    const setFormValue = <Key extends keyof BulkAddChannelPayload>(name: Key, value: BulkAddChannelPayload[Key]) => {
        setFormValues((values: BulkAddChannelPayload) => ({
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

        const response = (await dispatch(bulkAddToChannel(formValues))) as BulkAddChannelEventResponse;
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
                {'Start'}
            </FormButton>
        </React.Fragment>
    );

    const form = (
        <ActualForm
            formValues={formValues}
            setFormValue={setFormValue}
        />
    );

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
                    <h1>{'Bulk add to ' + channelName}</h1>
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
    formValues: BulkAddChannelPayload;
    setFormValue: <Key extends keyof BulkAddChannelPayload>(name: Key, value: BulkAddChannelPayload[Key]) => Promise<{ error?: string }>;
}

const ActualForm = (props: ActualFormProps) => {
    const {formValues, setFormValue} = props;

    const components: FormComponentProps[] = [
        {
            label: 'File (.JSON format)',
            required: true,
            helpText: <div>
                <a
                    href='https://github.com/mattermost/mattermost-plugin-bulk-invite/blob/master/.readme/template.jsonc'
                    target='_blank'
                    rel='noreferrer'
                >{'Download a template'}</a> {'to ensure your file formatting is correct.'}</div>,
            element: (
                <input
                    id='bulk-add-channel-file'
                    onChange={(e) => {
                        if (e.target.files?.length === 1) {
                            setFormValue('file', e.target.files[0]);
                        }
                    }}
                    type='file'
                />
            ),
        },
        {
            label: 'Add existing members to the team if they don’t belong to it',
            required: false,
            disabledText: (
                <div>
                    {"You don't have permission to add users to this team."}
                </div>
            ),
            helpText: (
                <div>
                    {'Enabling this will add users from other teams to this one if they are present on the file.'}
                </div>
            ),
            element: (
                <input
                    id='bulk-add-channel-adde-to-team'
                    onChange={(e) => {
                        setFormValue('add_to_team', e.target.checked);
                    }}
                    value={String(formValues.add_to_team)}
                    // disabled={teamAddDisabled}
                    type='checkbox'
                />
            ),
        },
        {
            label: 'Add guests',
            required: false,
            helpText: (
                <div>
                    {'Add guests if they are present on the file. If this is unchecked guests wont be added to the team if the above setting is not checked.'}
                </div>
            ),
            element: (
                <input
                    id='bulk-add-channel-add-guests'
                    onChange={(e) => {
                        setFormValue('add_guests', e.target.checked);
                    }}
                    value={String(formValues.add_guests)}
                    type='checkbox'
                    checked={formValues.add_guests}
                />
            ),
        },
    ];

    return (
        <div className='bulk-add-channel-form'>
            {components.map((c) => (
                <FormComponent
                    {...c}
                    key={c.element.props.id}
                />
            ))}
        </div>
    );
};
