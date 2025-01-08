use anyhow::{Ok, Result};
use auth_handler::login_handler;
use clap::Parser;
use cli_struct::{App, Command, ConfigVariableOptions};
use config_handler::{change_config, set_config_name_as_default, IndividualConfig};
use env_handler::{pull_handler, push_handler};

mod api_util;
mod auth_handler;
mod cli_struct;
mod config_handler;
mod env_handler;

fn main() -> Result<()> {
    let cli = App::parse();

    match cli.command {
        Command::Pull(value) => pull_handler(value)?,
        Command::Push(value) => push_handler(value)?,
        Command::Config(config_subcommand) => match config_subcommand {
            ConfigVariableOptions::Default(value) => set_config_name_as_default(value.name)?,
            ConfigVariableOptions::Upsert(value) => change_config(IndividualConfig {
                org: value.org,
                url: value.url,
                access_token: None,
                device_code: None,
                refresh_token: None,
                client_id: value.client_id,
            })?,
        },
        Command::Login(value) => login_handler(value)?,
    }

    Ok(())
}
