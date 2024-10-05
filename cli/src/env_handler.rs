use crate::cli_struct::RepoActionEnvArgs;
use crate::config_handler::{get_default_org, get_default_url};
use anyhow::{anyhow, Context, Result};
use base64::prelude::*;
use reqwest::{header::CONTENT_TYPE, Client, Response, StatusCode};
use serde::de::DeserializeOwned;
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

fn get_org(value: RepoActionEnvArgs) -> Result<String> {
    value
        .org
        .or(Some(get_default_org()?))
        .ok_or_else(|| anyhow::anyhow!("Org parameter is missing"))
}

fn get_env_path(value: RepoActionEnvArgs) -> Result<String> {
    value
        .path
        .to_str()
        .ok_or_else(|| anyhow::anyhow!("Invalid env path"))
        .and_then(|path| Ok(path.to_owned()))
}

async fn treat_api_err<T: DeserializeOwned>(response: Response) -> Result<T> {
    if let Err(err) = response.error_for_status_ref() {
        let status = err.status();
        let text = response.text().await?;

        return Err(anyhow!(
            "Request failed with status {}: {}",
            status.unwrap_or(StatusCode::INTERNAL_SERVER_ERROR),
            text
        ));
    }

    Ok(response.json::<T>().await?)
}

#[tokio::main]
pub async fn pull_handler(value: RepoActionEnvArgs) -> Result<()> {
    let url = get_default_url()?;
    let org = get_org(value.clone())?;
    let path = get_env_path(value.clone())?;
    let request_url = format!(
        "{}/orgs/{}/repos/{}?env={}",
        url, org, value.repository, value.env
    );

    let result =
        treat_api_err::<PullResponse>(Client::new().get(&request_url).send().await?).await?;

    let env = BASE64_STANDARD
        .decode(&result.file)
        .expect("Failed to decode base64 data");
    let mut file = File::create(path)?;

    file.write_all(&env)?;

    Ok(())
}

#[tokio::main]
pub async fn push_handler(value: RepoActionEnvArgs) -> Result<()> {
    let path = get_env_path(value.clone())?;
    let url = get_default_url()?;
    let org = get_org(value.clone())?;
    let content = std::fs::read_to_string(path.clone())
        .with_context(|| format!("could not read file `{}`", path))?;
    let request_url = format!(
        "{}/orgs/{}/repos/{}?env={}",
        url, org, value.repository, value.env,
    );
    let request_body = json!({
        "b64String": BASE64_STANDARD.encode(content)
    });

    let result = treat_api_err::<PushResponse>(
        Client::new()
            .post(request_url)
            .header(CONTENT_TYPE, "application/json")
            .json(&request_body)
            .send()
            .await?,
    )
    .await?;

    println!("{}", result.message);

    Ok(())
}
