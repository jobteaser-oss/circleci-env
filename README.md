# CircleCI Environment Variables Management

## Abstract
When the amount of CircleCI project grows it's can be complicate to create, update
or update environment variables for each project. And sometime you would add
environment variable without frontend tracking. This CLI wrap the CircleCI API
to manage the environment variable.

## Usage
The CLI expose 4 commands: list, get, set and del.

    # list the environment variables in a project
    $> circleci-env \
        --token $CIRCLECI_TOKEN \
        --vcs-type github \
        --username jobteaser-oss \
        --project someproject \
        list

    # get a environment variable in a project
    $> circleci-env \
        --token $CIRCLECI_TOKEN \
        --vcs-type github \
        --username jobteaser-oss \
        --project someproject \
        get FOO

    # set (create or update) a environment variable in a project
    $> circleci-env \
        --token $CIRCLECI_TOKEN \
        --vcs-type github \
        --username jobteaser-oss \
        --project someproject \
        set FOO BAR

    # delete a environment variable in a project
    $> circleci-env \
        --token $CIRCLECI_TOKEN \
        --vcs-type github \
        --username jobteaser-oss \
        --project someproject \
        del FOO

The flag `--help` or `-h` can give you more information about each command.

## Build

Building the project requires go >= 1.12.x
You can build the service with:

    make build

The final binary is available in the bin directory.
