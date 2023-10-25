import React, {PureComponent} from 'react';
import PropTypes from 'prop-types';

interface FormButtonProps {
    executing?: boolean;
    executingMessage?: React.ReactNode;
    defaultMessage?: React.ReactNode;
    btnClass?: string;
    extraClasses?: string;

    id?: string;
    type?: 'button' | 'submit' | 'reset';
    disabled?: boolean;
    onClick?: React.MouseEventHandler<HTMLButtonElement>;
}

const FormButton = ({
    executing = false,
    disabled = false,
    executingMessage = 'Creating',
    defaultMessage = 'Create',
    btnClass = 'btn btn-primary',
    extraClasses = '',
    ...props
}: FormButtonProps) => {
    let contents: React.ReactNode;
    if (executing) {
        contents = (
            <span>
                <span
                    className='fa fa-spin fa-spinner'
                    title={'Loading Icon'}
                />
                {executingMessage}
            </span>
        );
    } else {
        contents = defaultMessage;
    }

    let className = 'save-button btn ' + btnClass;

    if (extraClasses) {
        className += ' ' + extraClasses;
    }

    return (
        <button
            id={props.id}
            type={props.type}
            className={className}
            disabled={disabled}
            onClick={props.onClick}
        >
            {contents}
        </button>
    );
};

export default FormButton;
