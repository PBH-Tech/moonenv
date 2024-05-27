use clap::{ Args, Parser, Subcommand, ValueEnum};
use anyhow::{Context, Result};
use serde::{Deserialize,Serialize};
use serde_json::json;
use base64::prelude::*;
use std::fmt;
use std::io::Write;
use std::fs::File;
use reqwest::{Client, header::CONTENT_TYPE};

/// Manages environment helping saving and pulling it
#[derive(Parser, Debug)]
pub struct App {
    #[clap(subcommand)]
    command: Command
}

#[derive(Deserialize, Debug)]
struct PushResponse {
    message: String
}

#[derive(Deserialize, Debug)]
struct PullResponse {
    file: String
}

#[derive(Subcommand, Debug)]
enum Command {
    /// Pulls the .env from the indicated repository
    Pull(RepoActionEnvArgs),

    /// Pushed the .env file located on the path where the command has been executed to the repository
    Push(RepoActionEnvArgs)
}

#[derive(Clone, ValueEnum, Debug, Serialize)]
enum Environment {
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
struct RepoActionEnvArgs {
    /// The organization that owns the repository. Ensure that you have the necessary access permissions.
    org: String,

    /// The specific repository within the given organization where the `.env` file is located.
    repository: String,

    /// Environment where to find the .env file
    env: Environment
}



fn main() -> Result<()> {
    let cli = App::parse();

    match cli.command {
        Command::Pull(value) => { pull_handler(value)? }
        Command::Push(value) => { push_handler(value)? }
    }


    Ok(())
}

#[tokio::main]
async fn pull_handler(value: RepoActionEnvArgs)-> Result<()>{
    let path = "./.env"; // TODO: Turn path as an optional field
    let request_url = format!("https://t5m17jo2d8.execute-api.ap-southeast-2.amazonaws.com/dev/sendPullEnv?org={}&repo={}&env={}", value.org, value.repository, value.env);
    let response = Client::new().get(request_url).send().await?;

    response.error_for_status_ref()?;

    let result: PullResponse = response.json().await?;
    let env = BASE64_STANDARD.decode(result.file).expect("Failed to decode base64 data");
    let mut file = File::create(path)?;

    file.write_all(&env)?;
    
    Ok(())
}

#[tokio::main]
async fn push_handler(value: RepoActionEnvArgs) -> Result<()> {
    let path = "./.env"; // TODO: Turn path as an optional field
    let content = std::fs::read_to_string(path).with_context(|| format!("could not read file `{}`", path))?;
    let request_url = "https://t5m17jo2d8.execute-api.ap-southeast-2.amazonaws.com/dev/sendPushEnv"; // TODO: Convert to environment variable
    let request_body = json!({
        "org": value.org,
        "repo": value.repository,
        "env": value.env,
        "b64String": BASE64_STANDARD.encode(content)
    });
    let response = Client::new()
    .post(request_url)
    .header(CONTENT_TYPE, "application/json")
    .json(&request_body)
    .send().await?;

    response.error_for_status_ref()?;
    
    let result: PushResponse = response.json().await?;
    
    println!("{}", result.message);

    Ok(())
}
