use anyhow::{Ok, Result};
use clap::Parser;
use cli_struct::{App, Command, ConfigVariableOptions};
use config_handler::{change_config, set_default, IndividualConfig};
use env_handler::{pull_handler, push_handler};

mod cli_struct;
mod config_handler;
mod env_handler;

fn main() -> Result<()> {
    let cli = App::parse();

    match cli.command {
        Command::Pull(value) => pull_handler(value)?,
        Command::Push(value) => push_handler(value)?,
        Command::Config(config_subcommand) => match config_subcommand {
            ConfigVariableOptions::Default(value) => set_default(value.name)?,
            ConfigVariableOptions::Upsert(value) => change_config(IndividualConfig {
                org: value.name,
                url: value.url,
            })?,
        },
    }

    Ok(())
}
