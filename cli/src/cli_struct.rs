use clap::{Args, Parser, Subcommand, ValueEnum};
use serde::Serialize;
use std::fmt;

/// Manages environment helping saving and pulling it
#[derive(Parser, Debug)]
pub struct App {
    #[clap(subcommand)]
    pub command: Command,
}

#[derive(Subcommand, Debug)]
pub enum Command {
    /// Pulls the .env from the indicated repository
    Pull(RepoActionEnvArgs),

    /// Pushed the .env file located on the path where the command has been executed to the repository
    Push(RepoActionEnvArgs),

    /// Changes the application's configuration settings.
    Config(ConfigVariable),
}

#[derive(Clone, ValueEnum, Debug, Serialize)]
pub enum Environment {
    Dev,
    Qa,
    Prod,
}

impl fmt::Display for Environment {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            Environment::Dev => write!(f, "Dev"),
            Environment::Qa => write!(f, "Qa"),
            Environment::Prod => write!(f, "Prod"),
        }
    }
}

#[derive(Args, Debug)]
pub struct RepoActionEnvArgs {
    /// The organization that owns the repository. Ensure that you have the necessary access permissions.
    pub org: String,

    /// The specific repository within the given organization where the `.env` file is located.
    pub repository: String,

    /// Environment where to find the .env file
    pub env: Environment,
}

#[derive(Args, Debug)]
pub struct ConfigVariable {
    /// A friendly name for identifying this configuration.
    pub name: String,

    #[clap(short, long)]
    /// The full URL to the server.
    /// If provided, it should be a valid URL format, e.g., "https://example.com".
    pub url: Option<String>,
}
