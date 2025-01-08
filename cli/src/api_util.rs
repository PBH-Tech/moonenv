use anyhow::{anyhow, Result};
use reqwest::{Response, StatusCode};
use serde::de::DeserializeOwned;

pub async fn treat_api_err<T: DeserializeOwned>(response: Response) -> Result<T> {
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
