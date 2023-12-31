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
    file?: File
    users: string[];
    channel_id: string;
}

export default function BulkAddChannelForm(props: Props) {
    const dispatch = useDispatch();
    const modalProps = useSelector(getBulkAddChannelModal);
    const [storedError, setStoredError] = useState('');
    const [submitting, setSubmitting] = useState(false);
    const [channelName, setChannelName] = useState('');
    const [formValues, setFormValues] = useState<BulkAddChannelPayload>({
        add_to_team: false,
        users: [],
        channel_id: modalProps?.channelId || '',
    });

    const loadChannelInfo = useCallback(async (channelId: string): Promise<Channel | null | undefined> => {
        const response = await getChannelInfo(channelId);

        if (response.error) {
            return null;
        }

        return response.channel;
    }, []);

    if (!channelName && modalProps?.channelId) {
        loadChannelInfo(modalProps.channelId).then((channel) => {
            if (channel) {
                setChannelName(channel.display_name);
            }
        });
    }

    const setFormValue = <Key extends keyof BulkAddChannelPayload>(name: Key, value: BulkAddChannelPayload[Key]) => {
        setFormValues((values: BulkAddChannelPayload) => ({
            ...values,
            [name]: value,
        }));
    };

    const handleClose = (e?: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
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

        const response = await bulkAddToChannel(formValues);
        if (response.error) {
            handleError(response.error);
            return;
        }

        handleClose();
    };

    const footer = (
        <React.Fragment>
            <FormButton
                defaultMessage='Cancel'
                btnClass='btn-tertiary'
                onClick={handleClose}
            />
            <FormButton
                id='submit-button'
                type='submit'
                executing={submitting}
                executingMessage='Adding users'
                defaultMessage='Add users'
            />
        </React.Fragment>
    );

    const form = (
        <ActualForm
            formValues={formValues}
            setFormValue={setFormValue}
        />
    );

    let error: React.ReactNode = null;
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
            <Modal.Body>
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
    setFormValue: <Key extends keyof BulkAddChannelPayload>(name: Key, value: BulkAddChannelPayload[Key]) => void;
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
                    accept='.json'
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
