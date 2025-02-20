use std::{
    sync::Arc,
    thread::{sleep, spawn},
    time::{Duration, SystemTime, UNIX_EPOCH},
};

use crate::moonenv_config::MoonenvConfig;
use crate::{api_util::treat_api_err, cli_struct::OrgActionAuthArgs};
use anyhow::{anyhow, Ok, Result};
use reqwest::{Client, StatusCode};
use serde::Deserialize;

#[derive(Deserialize, Debug)]
struct OAuthSetOfTokenResult {
    #[serde(alias = "authorizationUri")]
    authorization_uri: String,

    #[serde(alias = "deviceCode")]
    device_code: String,
}

#[derive(Deserialize, Debug, Clone)]
struct OAuthTokenResult {
    #[serde(alias = "idToken")]
    id_token: String,

    #[serde(alias = "refreshToken")]
    refresh_token: String,

    #[serde(alias = "expiresIn")]
    expires_in: u16,
}

#[derive(Deserialize, Debug, Clone)]
struct OAuthRefreshTokenResult {
    #[serde(alias = "idToken")]
    id_token: String,

    #[serde(alias = "expiresIn")]
    expires_in: u16,
}

#[tokio::main]
pub async fn login_handler(value: OrgActionAuthArgs) -> Result<()> {
    let mut moonenv_config = MoonenvConfig::new();
    let org = Arc::new(moonenv_config.get_org(value.org)?);
    let url = moonenv_config.get_url(&org)?;
    let client_id = moonenv_config.get_client_id(&org)?;
    let uri = format!("{}/auth/token?client_id={}", url, client_id);
    let set_of_token_result =
        treat_api_err::<OAuthSetOfTokenResult>(Client::new().get(&uri).send().await?).await?;
    open::that(format!("https://{}", set_of_token_result.authorization_uri))?;
    let org_clone = Arc::clone(&org);
    let (set_of_token_result, login_result) =
        spawn(move || fetch_login_result(set_of_token_result, &org_clone))
            .join()
            .map_err(|e| anyhow::Error::msg(format!("Login failed: {:?}", e)))??;
    let mut config = moonenv_config.get_config(&org)?;

    config.access_token = Some(login_result.id_token); // TODO: weird, but access token is ID Token
    config.device_code = Some(set_of_token_result.device_code);
    config.refresh_token = Some(login_result.refresh_token);
    config.access_token_expires_at = Some(get_expires_at(login_result.expires_in)?);

    let _ = moonenv_config.change_config(config);

    Ok(())
}

#[tokio::main]
async fn fetch_login_result(
    set_of_token_result: OAuthSetOfTokenResult,
    org: &str,
) -> Result<(OAuthSetOfTokenResult, OAuthTokenResult)> {
    let mut moonenv_config = MoonenvConfig::new();
    let url = moonenv_config.get_url(org)?;
    let client_id = moonenv_config.get_client_id(org)?;
    let uri = format!(
        "{}/auth/token?client_id={}&device_code={}&grant_type=urn:ietf:params:oauth:grant-type:device_code",
        url, client_id, set_of_token_result.device_code
    );
    let token_result: OAuthTokenResult;

    loop {
        sleep(Duration::from_millis(3100));

        let result = Client::new().get(&uri).send().await?;
        let status = result.status();

        if status.is_success() {
            token_result = result.json::<OAuthTokenResult>().await?;
            break;
        } else if status.as_u16() == StatusCode::GONE {
            return Err(anyhow!("Session is expired. Try to login again!"));
        }
    }

    Ok((set_of_token_result, token_result))
}

pub async fn get_access_token(org: &str) -> Result<String> {
    let mut moonenv_config = MoonenvConfig::new();
    let config = moonenv_config.get_config(org)?;
    let mut access_token = config.access_token;
    let now = get_duration_since_unix_epoch();
    let expires_at = config.access_token_expires_at.unwrap_or(now);

    if Option::is_none(&access_token) || expires_at < now {
        access_token = Some(refresh_token(org).await?);
    }

    Ok(access_token.ok_or(anyhow::anyhow!("Access token is not defined"))?)
}

fn get_expires_at(expires_in: u16) -> Result<Duration> {
    let now = get_duration_since_unix_epoch();

    Ok(now + Duration::from_secs(expires_in.into()))
}

fn get_duration_since_unix_epoch() -> Duration {
    SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .expect("Time went backwards")
}

async fn refresh_token(org: &str) -> Result<String> {
    let mut moonenv_config = MoonenvConfig::new();
    let mut config = moonenv_config.get_config(org)?;
    let refresh_token = config.refresh_token.clone().ok_or(anyhow::anyhow!(
        "No refresh token found. Try to login first"
    ))?;
    let device_code = config
        .device_code
        .clone()
        .ok_or(anyhow::anyhow!("No device code found. Try to login first."))?;
    let url = moonenv_config.get_url(org)?;
    let uri = format!("{}/auth/refresh-token?device_code={}", url, device_code);
    let result = treat_api_err::<OAuthRefreshTokenResult>(
        Client::new()
            .post(uri)
            .bearer_auth(refresh_token)
            .send()
            .await?,
    )
    .await?;

    config.access_token = Some(result.id_token.clone());
    config.access_token_expires_at = Some(get_expires_at(result.expires_in)?);

    let _ = moonenv_config.change_config(config);

    Ok(result.id_token)
}
