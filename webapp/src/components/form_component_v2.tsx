import React from 'react';

export type Props = {
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
        label,
        required,
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
            <div className='help-text'>
                {helpText}
            </div>
            <div>
                {(element.props.type != "checkbox") && element}
            </div>
        </div>
    );
}
