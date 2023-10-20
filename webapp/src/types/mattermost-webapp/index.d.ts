import React from "react"

export interface PluginRegistry {
    registerReducer(reducer: Reducer<any, any>)
    registerChannelHeaderMenuAction(title: string, action: ChannelHeaderMenuAction, shouldRender?: () => boolean)
    registerRootComponent(component: React.ComponentType<any>)
}
