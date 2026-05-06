/* eslint-disable @typescript-eslint/no-explicit-any */
import React from 'react';
import {Reducer} from 'redux';

export interface PluginRegistry {
    registerReducer(reducer: Reducer<any, any>): void;
    registerChannelHeaderMenuAction(title: string, action: ChannelHeaderMenuAction, shouldRender?: () => boolean): void;
    registerRootComponent(component: React.ComponentType<any>): void;
}
