# Bulk Invite Mattermost Plugin

[![Delivery status](https://github.com/mattermost/mattermost-plugin-bulk-invite/actions/workflows/cd.yml/badge.svg)](https://github.com/mattermost/mattermost-plugin-bulk-invite/actions/workflows/cd.yml)
[![Code Coverage](https://img.shields.io/codecov/c/github/mattermost/mattermost-plugin-bulk-invite/master)](https://codecov.io/gh/mattermost/mattermost-plugin-bulk-invite)
[![Release](https://img.shields.io/github/v/release/mattermost/mattermost-plugin-bulk-invite)](https://github.com/mattermost/mattermost-plugin-bulk-invite/releases/latest)
[![HW](https://img.shields.io/github/issues/mattermost/mattermost-plugin-bulk-invite/Up%20For%20Grabs?color=dark%20green&label=Help%20Wanted)](https://github.com/mattermost/mattermost-plugin-bulk-invite/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22Up+For+Grabs%22+label%3A%22Help+Wanted%22)

This plugin allows you to add users to a channel in bulk by uploading a JSON file.

## Features

- Allows adding users to a channel in bulk by uploading a JSON file.
    - Supports using `user_id` and `username`.
- (Optionally) Adds the users to the team if they don't belong to it.
- (Optionally) Invite guest users too, if provied.

## Installation

1. Clone this repository.
2. Build and upload the plugin manually:
    1. Run `make dist` to build the plugin.
    2. Go to **System Console > Plugins > Management** in your Mattermost instance
    3. Upload the plugin located in the `dist/` folder.
3. (or) Upload the plugin directly with a command:
    1. Set the environment variables:
        -  `MM_SERVICESETTINGS_SITEURL` to your Mattermost URL.
        -  `MM_ADMIN_USERNAME` to your Mattermost username.
        -  `MM_ADMIN_PASSWORD` to your Mattermost password.
    2. Run `make deploy` to build and upload the plugin.

## Usage

After successful installation:

1. Craft a JSON file following the [following format](./.readme/template.jsonc).
2. Launch the plugin from the channel header or channel intro:
    - **Channel name > Bulk Invite**

        ![Channel header](./.readme/channel-header-menu.png)
    - **Channel intro > Bulk Invite** (only visible to channel admins)

        ![Channel intro](./.readme/channel-intro-button.png)

3. You will be presented with a modal where you can upload the JSON file:

    ![Bulk invite modal](./.readme/bulk-invite-modal.png)

    - **File**: Upload a JSON file following the [following format](./.readme/template.jsonc).
    - **Invite members to the team**: If checked, the users will be added to the team if they are not already members. Otherwise they will be skipped.
    - **Invite guests**: If checked, guest users on the list will be added to the channel (and team if the above is checked). Otherwise they will be skipped.

4. The plugin will display it's progress in the channel:

    ![Bulk invite progress](./.readme/result-channel.png)

    ![Bulk invite progress](./.readme/result-channel-thread.png)


## Contribute

If you would like to help improve this plugin, feel free to submit a pull request.
You can also check for open issues and see if there's anything you can help with.
