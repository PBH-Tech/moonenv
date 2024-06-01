use crate::cli_struct::RepoActionEnvArgs;
use anyhow::{Context, Result};
use base64::prelude::*;
use reqwest::{header::CONTENT_TYPE, Client};
use serde::Deserialize;
use serde_json::json;
use std::fs::File;
use std::io::Write;

#[derive(Deserialize, Debug)]
pub struct PushResponse {
    pub message: String,
}

#[derive(Deserialize, Debug)]
pub struct PullResponse {
    pub file: String,
}

#[tokio::main]
pub async fn pull_handler(value: RepoActionEnvArgs) -> Result<()> {
    let path = "./.env"; // TODO: Turn path as an optional field
    let request_url = format!("https://t5m17jo2d8.execute-api.ap-southeast-2.amazonaws.com/dev/sendPullEnv?org={}&repo={}&env={}", value.org, value.repository, value.env);
    let response = Client::new().get(request_url).send().await?;

    response.error_for_status_ref()?;

    let result: PullResponse = response.json().await?;
    let env = BASE64_STANDARD
        .decode(result.file)
        .expect("Failed to decode base64 data");
    let mut file = File::create(path)?;

    file.write_all(&env)?;

    Ok(())
}

#[tokio::main]
pub async fn push_handler(value: RepoActionEnvArgs) -> Result<()> {
    let path = "./.env"; // TODO: Turn path as an optional field
    let content =
        std::fs::read_to_string(path).with_context(|| format!("could not read file `{}`", path))?;
    let request_url = "https://t5m17jo2d8.execute-api.ap-southeast-2.amazonaws.com/dev/sendPushEnv"; // TODO: Convert it to environment variable
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
        .send()
        .await?;

    response.error_for_status_ref()?;

    let result: PushResponse = response.json().await?;

    println!("{}", result.message);

    Ok(())
}
