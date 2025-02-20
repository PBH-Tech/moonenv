use crate::api_util::treat_api_err;
use crate::auth_handler::get_access_token;
use crate::cli_struct::{ConfigVariableOptions, RepoActionEnvArgs};
use crate::moonenv_config::{IndividualConfig, MoonenvConfig};
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

impl RepoActionEnvArgs {
    pub fn new(args: RepoActionEnvArgs) -> Self {
        Self {
            org: args.org,
            env: args.env,
            path: args.path,
            repository: args.repository,
        }
    }

    #[tokio::main]
    pub async fn pull_handler(&mut self) -> Result<()> {
        let mut moonenv_config = MoonenvConfig::new();
        let path = self.get_env_path()?;
        let org = moonenv_config.get_org(self.org.clone())?;
        let request_url = self.get_request_url()?;

        let result = treat_api_err::<PullResponse>(
            Client::new()
                .get(&request_url)
                .bearer_auth(get_access_token(org.borrow()).await?)
                .send()
                .await?,
        )
        .await?;

        let env = BASE64_STANDARD
            .decode(result.file)
            .expect("Failed to decode base64 data");
        let mut file = File::create(path)?;

        file.write_all(&env)?;

        Ok(())
    }

    #[tokio::main]
    pub async fn push_handler(&mut self) -> Result<()> {
        let mut moonenv_config = MoonenvConfig::new();
        let path = self.get_env_path()?;
        let content = std::fs::read_to_string(path.clone())
            .with_context(|| format!("Could not read file `{}`", path))?;
        let org = moonenv_config.get_org(self.org.clone())?;
        let request_url = self.get_request_url()?;
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

    fn get_env_path(&mut self) -> Result<String> {
        self.path
            .to_str()
            .ok_or_else(|| anyhow::anyhow!("Invalid env path"))
            .map(|path| path.to_owned())
    }

    fn get_request_url(&mut self) -> Result<String> {
        let mut moonenv_config = MoonenvConfig::new();
        let org = moonenv_config.get_org(self.org.clone())?;
        let url = moonenv_config.get_url(org.borrow())?;

        Ok(format!(
            "{}/orgs/{}/repos/{}?env={}",
            url, org, self.repository, self.env
        ))
    }
}

impl ConfigVariableOptions {
    pub fn execute(option: ConfigVariableOptions) -> Result<()> {
        match option {
            ConfigVariableOptions::Default(value) => Self::set_config_name_as_default(value.name),
            ConfigVariableOptions::Upsert(value) => Self::change_config(IndividualConfig {
                org: value.org,
                url: value.url,
                access_token: None,
                access_token_expires_at: None,
                device_code: None,
                refresh_token: None,
                client_id: value.client_id,
            }),
        }
    }

    fn set_config_name_as_default(name: String) -> Result<()> {
        let mut moonenv_config = MoonenvConfig::new();
        moonenv_config.set_config_name_as_default(name)?;

        Ok(())
    }

    fn change_config(config: IndividualConfig) -> Result<()> {
        let mut moonenv_config = MoonenvConfig::new();
        moonenv_config.change_config(config)?;

        Ok(())
    }
}
