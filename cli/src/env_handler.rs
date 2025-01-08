use crate::api_util::treat_api_err;
use crate::auth_handler::get_access_token;
use crate::cli_struct::RepoActionEnvArgs;
use crate::config_handler::{get_org, get_url};
use anyhow::{Context, Result};
use base64::prelude::*;
use reqwest::{header::CONTENT_TYPE, Client};
use serde::Deserialize;
use serde_json::json;
use std::borrow::Borrow;
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

fn get_env_path(value: RepoActionEnvArgs) -> Result<String> {
    value
        .path
        .to_str()
        .ok_or_else(|| anyhow::anyhow!("Invalid env path"))
        .and_then(|path| Ok(path.to_owned()))
}

fn get_request_url(value: RepoActionEnvArgs) -> Result<String> {
    // If no org, the default one is used
    let org = get_org(value.org)?;
    let url = get_url(org.borrow())?;

    Ok(format!(
        "{}/orgs/{}/repos/{}?env={}",
        url, org, value.repository, value.env
    ))
}

#[tokio::main]
pub async fn pull_handler(value: RepoActionEnvArgs) -> Result<()> {
    let path = get_env_path(value.clone())?;
    let org = get_org(value.org.clone())?;
    let request_url = get_request_url(value)?;

    let result = treat_api_err::<PullResponse>(
        Client::new()
            .get(&request_url)
            .bearer_auth(get_access_token(org.borrow()).await?)
            .send()
            .await?,
    )
    .await?;

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
    let content = std::fs::read_to_string(path.clone())
        .with_context(|| format!("Could not read file `{}`", path))?;
    let org = get_org(value.org.clone())?;
    let request_url = get_request_url(value)?;
    let request_body = json!({
        "b64String": BASE64_STANDARD.encode(content)
    });

    let result = treat_api_err::<PushResponse>(
        Client::new()
            .post(request_url)
            .bearer_auth(get_access_token(org.borrow()).await?)
            .header(CONTENT_TYPE, "application/json")
            .json(&request_body)
            .send()
            .await?,
    )
    .await?;

    println!("{}", result.message);

    Ok(())
}
