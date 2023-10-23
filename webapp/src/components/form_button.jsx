import React from 'react';
import PropTypes from 'prop-types';

const FormButton = ({
    saving,
    disabled,
    savingMessage,
    defaultMessage,
    btnClass,
    extraClasses,
    ...props
}) => {
    let contents;
    if (saving) {
        contents = (
            <span>
                <span
                    className='fa fa-spin fa-spinner'
                    title={'Loading Icon'}
                />
                {savingMessage}
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
            id='saveSetting'
            className={className}
            disabled={disabled}
            {...props}
        >
            {contents}
        </button>
    );
};

FormButton.propTypes = {
    executing: PropTypes.bool,
    disabled: PropTypes.bool,
    executingMessage: PropTypes.node,
    defaultMessage: PropTypes.node,
    btnClass: PropTypes.string,
    extraClasses: PropTypes.string,
    saving: PropTypes.bool,
    savingMessage: PropTypes.string,
    type: PropTypes.string,
};

FormButton.defaultProps = {
    disabled: false,
    savingMessage: 'Creating',
    defaultMessage: 'Create',
    btnClass: 'btn-primary',
    extraClasses: '',
};

export default FormButton;
