import React from 'react';

export type Props = {
    inputId?: string;
    label: React.ReactNode;
    element: React.ReactElement;
    helpText?: JSX.Element;
    required?: boolean;
    hideRequiredStar?: boolean;
    type?: string;
}

export function FormComponentV2(props: Props) {
    const {
        element,
        helpText,
        inputId,
        label,
        required,
        hideRequiredStar,
    } = props;

    return (
        <div className='form-group'>
            {(element.props.type == "checkbox") && element}
            {label &&
                <label
                    className='control-label margin-bottom x2'
                    htmlFor={inputId}
                >
                    {label}
                </label>
            }
            {required && !hideRequiredStar &&
                <span
                    className='error-text'
                    style={{marginLeft: '3px'}}
                >
                    {'*'}
                </span>
            }
            <div className='help-text'>
                {helpText}
            </div>
            <div>
                {(element.props.type != "checkbox") && element}
            </div>
        </div>
    );
}
