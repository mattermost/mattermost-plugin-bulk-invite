export interface PluginRegistry {
    // Add more if needed from https://developers.mattermost.com/extend/plugins/webapp/reference
    registerPostDropdownMenuAction(text: string, action: (postId: string) => void, filter: (postId: string) => bool)
}
