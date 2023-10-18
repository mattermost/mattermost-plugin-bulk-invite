import React from 'react';
import {useSelector, useDispatch} from 'react-redux';

import {Modal} from 'react-bootstrap';

import {isBulkInviteChannelModalVisible} from '@/selectors';

import BulkInviteChannelForm from '../forms/bulk_invite_channel_form';
import { closeBulkInviteChannelModal } from '@/actions';

type Props = {};

import './bulk_invite_modal.scss';

export default function BulkInviteChannelModal(props: Props) {
    const visible = useSelector(isBulkInviteChannelModalVisible);

    const dispatch = useDispatch();
    const close = () => dispatch(closeBulkInviteChannelModal());

    if (!visible) {
        return null;
    }

    const content = (
        <BulkInviteChannelForm
            {...props}
            close={close}
        />
    );

    return (
        <Modal
            id='bulk-invite-channel-modal'
            dialogClassName='a11y__modal channel-invite'
            show={visible}
            onHide={close}
            onExited={close}
        >
            <Modal.Header closeButton={true}></Modal.Header>
            {content}
        </Modal>
    );
}
