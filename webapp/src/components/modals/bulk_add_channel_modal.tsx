import React from 'react';
import {useSelector, useDispatch} from 'react-redux';

import {Modal} from 'react-bootstrap';

import {isBulkAddChannelModalVisible} from '@/selectors';

import BulkAddChannelForm from '../forms/bulk_add_channel_form';
import {closeBulkAddChannelModal} from '@/actions';

type Props = {};

import './bulk_add_channel_modal.scss';

export default function BulkAddChannelModal(props: Props) {
    const visible = useSelector(isBulkAddChannelModalVisible);

    const dispatch = useDispatch();
    const close = () => dispatch(closeBulkAddChannelModal());

    if (!visible) {
        return null;
    }

    const content = (
        <BulkAddChannelForm
            {...props}
            close={close}
        />
    );

    return (
        <Modal
            id='bulk-add-channel-modal'
            dialogClassName='a11y__modal channel-invite'
            show={visible}
            onHide={close}
            onExited={close}
        >
            <Modal.Header closeButton={true}/>
            {content}
        </Modal>
    );
}
