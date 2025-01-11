# Moonenv 🌚

Moonenv is a simple CLI that helps you manage, version, and share environment variables.

Inspired by [dotenv.org](https://www.dotenv.org/), Moonenv aims to provide a low-cost tool for self-management.

# Available commands
The list below maps all the available commands, but when the tool is installed, the `--help` command is available to give more details for each command.
- pull: Pulls the environment variable file;
- push: Pushed the file and keep the historical changes;
- config: Changes the application's configuration settings;
- login: Initiates the server authentication process;
- help: Explain more about the commands and their subcommands;

# How does it work?

Moonenv is split into two parts: CLI and server. 
- The CLI is where you are going to execute the commands;
- The server is responsible for handling the command requests and providing the response for each of them: authenticating, saving (keeping the history of change) and providing the environment variables;


## CLI config
The CLI saves all the configurations in the `~/.moonenv/config` file. Its structure looks like this:
- default: The default profile to be used when the `-org` subcommand is not used;
- profiles: List of profiles set;
  - url: The URL for the server. It has to start with `https://`
  - org: The profile name that can be used at the `default` or `--org` subcommand.
  - device_code: When authenticated, the server generates a device code used during the refresh token life cycle.
  - access_token: Token to authenticate with the server.
  - access_token_expires_at: Access token expiration time.
  - refresh_token: Token to access a new access token when it is expired.
  - client_id: The OAuth client ID for the CLI tool.

# Running the server

IN PROGRESS

# Contributing

Refer to our contribution guidelines and [Code of Conduct for contributors](https://github.com/PBH-Tech/moonenv/blob/main/CODE_OF_CONDUCT.md).