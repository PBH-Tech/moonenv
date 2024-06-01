use anyhow::Result;
use clap::Parser;
use cli_struct::{App, Command};
use env_handler::{pull_handler, push_handler};

mod cli_struct;
mod env_handler;

fn main() -> Result<()> {
    let cli = App::parse();

    match cli.command {
        Command::Pull(value) => pull_handler(value)?,
        Command::Push(value) => push_handler(value)?,
    }

    Ok(())
}
