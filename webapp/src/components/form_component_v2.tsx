import React from 'react';

export type Props = {
    label: React.ReactNode;
    element: React.ReactElement;
    helpText?: JSX.Element;
    required?: boolean;
    hideRequiredStar?: boolean;
    disabledText?: JSX.Element;
    type?: string;
}

export function FormComponentV2(props: Props) {
    const {
        element,
        helpText,
        label,
        required,
        disabledText,
        hideRequiredStar,
    } = props;

    return (
        <div className='form-group'>
            <label
                className='control-label margin-bottom x2'
                htmlFor={element.props.id}
            >
                {(element.props.type == "checkbox") && element}
                {label}
            {required && !hideRequiredStar &&
                <span
                className='error-text'
                style={{marginLeft: '3px'}}
                >
                    {'*'}
                </span>
            }
            </label>
            {helpText && !element.props.disabled &&
                <div className='help-text'>
                    {helpText}
                </div>}
            {element.props.disabled && disabledText &&
                <div className='help-text disabled-text'>
                    {disabledText}
                </div>}
            <div>
                {(element.props.type != "checkbox") && element}
            </div>
        </div>
    );
}
