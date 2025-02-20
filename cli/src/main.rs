use anyhow::{Ok, Result};
use auth_handler::login_handler;
use clap::Parser;
use cli_struct::{App, Command, ConfigFileOptions, RepoActionEnvArgs};

mod api_util;
mod auth_handler;
mod cli_struct;
mod env_handler;
mod moonenv_config;

fn main() -> Result<()> {
    let cli = App::parse();

    match cli.command {
        Command::Pull(value) => RepoActionEnvArgs::new(value).pull_handler()?,
        Command::Push(value) => RepoActionEnvArgs::new(value).push_handler()?,
        Command::Config(config_subcommand) => ConfigFileOptions::execute(config_subcommand)?,
        Command::Login(value) => login_handler(value)?,
    }

    Ok(())
}
