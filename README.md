% duocli 8

# NAME

duocli - CLI Interface to Duo Admin API

# SYNOPSIS

duocli

```
[--config|-c]=[value]
[--help|-h]
[--version|-v]
```

**Usage**:

```
duocli [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--config, -c**="": load configuration from JSON `FILE`

**--help, -h**: show help

**--version, -v**: print the version


# COMMANDS

## user

manage users

### create

create a user

**--email, -e**="": email address of user

**--firstName, -f**="": first name of user

**--group, -g**="": add user to group, can be specified multiple times to add user to multiple groups

**--lastName, -l**="": last name of user

**--status, -s**="": status of user: active, disabled, or bypass (default: active)

**--username, -u**="": username

### get

get one or more users and display as JSON

**--username, -u**="": username, can be specified multiple times

### modify

modify a user's attributes, add or remove group membership

**--addgroup, -g**="": add user to groups, adds to existing memberships, and can be specified multiple times to add user to multiple groups

**--create, -c**: create user if not found

**--delgroup, -G**="": remove user from groups, removes from existing memberships, and can be specified multiple times to remove user from multiple groups

**--email, -e**="": email address of user

**--firstName, -f**="": first name of user

**--lastName, -l**="": last name of user

**--status, -s**="": status of user: active, disabled, or bypass

**--username, -u**="": username

### delete

delete user and any attached phones

**--phone, -P**: delete any phones found attached to the user before deleting the user

**--username, -u**="": username, can be specified multiple times
