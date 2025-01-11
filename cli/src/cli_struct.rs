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

    /// Pushes the .env file located on the path where the command has been executed to the repository
    Push(RepoActionEnvArgs),

    #[clap(subcommand)]
    /// Changes the application's configuration settings.
    Config(ConfigVariableOptions),

    /// Initiates a login process by redirecting to the login page on the browser
    Login(OrgActionAuthArgs),
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
            // TODO: change it to be dynamic; The client can set any environment
            Environment::Dev => write!(f, "Dev"),
            Environment::Qa => write!(f, "Qa"),
            Environment::Prod => write!(f, "Prod"),
        }
    }
}

#[derive(Args, Debug, Clone)]
pub struct RepoActionEnvArgs {
    #[clap(short, long)]
    /// The organization that owns the repository.
    /// Make sure you have the necessary access permissions.
    /// If unspecified, the organization name is taken from the default configuration profile.
    pub org: Option<String>,

    #[clap(short, long, default_value = "./.env")]
    /// Path to the environment variable file.
    pub path: std::path::PathBuf,

    /// The specific repository within the given organization where the `.env` file is located.
    pub repository: String,

    /// Environment where to find the .env file.
    pub env: Environment,
}

#[derive(Args, Debug, Clone)]
pub struct OrgActionAuthArgs {
    #[clap(short, long)]
    /// The organization that you want to login or logout.
    /// If unspecified, the organization name is taken from the default configuration profile.
    pub org: Option<String>,
}

#[derive(Subcommand, Debug)]
pub enum ConfigVariableOptions {
    /// Creates or updates a profile with specified settings.
    /// Use this command to either create a new profile or update an existing one with new values.
    Upsert(ConfigVariableUpsert),

    /// Sets the currently selected profile as the default for the application.
    /// This command will modify the application's settings so that the specified profile
    Default(ConfigVariableChangeDefault),
}

#[derive(Args, Debug)]
pub struct ConfigVariableUpsert {
    /// The org identifier for this configuration.
    pub org: String,

    #[clap(short, long)]
    /// The full URL to the server.
    /// If provided, it should be a valid URL format, e.g., <https://example.com>.
    pub url: Option<String>,

    #[clap(short, long)]
    /// The CLI client ID for authorization.
    pub client_id: Option<String>,
}

#[derive(Args, Debug)]
pub struct ConfigVariableChangeDefault {
    /// The name of the profile to set as default.
    pub name: String,
}
