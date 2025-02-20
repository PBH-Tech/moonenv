use anyhow::{Ok, Result};
use auth_handler::login_handler;
use clap::Parser;
use cli_struct::{App, Command, ConfigVariableOptions};
use env_handler::{pull_handler, push_handler};
use moonenv_config::{IndividualConfig, MoonenvConfig};

mod api_util;
mod auth_handler;
mod cli_struct;
mod env_handler;
mod moonenv_config;

fn main() -> Result<()> {
    let cli = App::parse();
    let mut moonenv_config = MoonenvConfig::new();

    match cli.command {
        Command::Pull(value) => pull_handler(value)?,
        Command::Push(value) => push_handler(value)?,
        Command::Config(config_subcommand) => match config_subcommand {
            ConfigVariableOptions::Default(value) => {
                moonenv_config.set_config_name_as_default(value.name)?
            }
            ConfigVariableOptions::Upsert(value) => {
                moonenv_config.change_config(IndividualConfig {
                    org: value.org,
                    url: value.url,
                    access_token: None,
                    access_token_expires_at: None,
                    device_code: None,
                    refresh_token: None,
                    client_id: value.client_id,
                })?
            }
        },
        Command::Login(value) => login_handler(value)?,
    }

    Ok(())
}
