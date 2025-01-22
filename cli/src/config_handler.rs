//! Config Handler
//!
//! This module is responsible for saving and reading the config in the ~/.moonenv/config file.
//! It contains functions that validates and gives the configuration.

use anyhow::{anyhow, Ok, Result};
use confy::{self};
use dirs::home_dir;
use serde::{Deserialize, Serialize};
use std::{borrow::Borrow, path::PathBuf, time::Duration};

/// The ConfigHandler
trait ConfigHandler {
    /// Returns the file path for the configuration file
    fn get_default_config_file_path() -> Result<PathBuf> {
        let mut home = home_dir().ok_or_else(|| anyhow!("HOME environment variable not set"))?;

        home.push(".moonenv");
        home.push("config");

        Ok(home.clone())
    }

    /// Saves the configuration
    fn save_config(moonenv_config: MoonenvConfig, path: PathBuf) -> Result<()> {
        let _ = confy::store_path(path, moonenv_config);

        Ok(())
    }

    /// Returns the configuration
    fn get_config(config_path: PathBuf) -> Result<MoonenvConfig> {
        confy::load_path(config_path).map_err(|e| anyhow!("Failed to load configuration: {}", e))
    }
}

/// The configuration struct, responsible to informs the default profile and all profiles that are available
#[derive(Debug, Serialize, Deserialize)]
pub struct MoonenvConfig {
    /// Default profile to be used if the "--org" flag is not used
    pub default: Option<String>,

    /// The list o available profiles.
    pub profiles: Vec<IndividualConfig>,
}

/// The individual profile configuration
#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct IndividualConfig {
    /// URL for the server API
    pub url: Option<String>,

    /// The profile name
    pub org: String,

    /// When authenticated, the server will return a device code that will be used to refresh the access token
    pub device_code: Option<String>,

    /// Access token used to access the server's endpoint
    pub access_token: Option<String>,

    /// The access token expiration time
    pub access_token_expires_at: Option<Duration>,

    /// Token to refresh the access token
    pub refresh_token: Option<String>,

    /// The client ID for authentication OAuth server
    pub client_id: Option<String>,
}

impl Default for IndividualConfig {
    fn default() -> Self {
        Self {
            org: "moonenv".to_string(),
            url: Some("www.moonenv.app".to_string()),
            access_token: None,
            access_token_expires_at: None,
            device_code: None,
            refresh_token: None,
            client_id: None,
        }
    }
}

impl ::std::default::Default for MoonenvConfig {
    fn default() -> Self {
        Self {
            default: Some("moonenv".into()),
            profiles: vec![IndividualConfig::default()],
        }
    }
}

impl ConfigHandler for MoonenvConfig {}

/// Change an individual configuration;
/// If it already exists, so it overwrites;
/// Otherwise, it pushes a new profile in the profile list;
/// And finally, it saves the config in the file
pub fn change_config(new_config: IndividualConfig) -> Result<()> {
    let default_path = MoonenvConfig::get_default_config_file_path()?;
    let mut moonenv_config = MoonenvConfig::get_config(default_path.clone())?;
    let config_org = new_config.org.clone();

    if let Some(individual_config) = moonenv_config
        .profiles
        .iter_mut()
        .find(|config| config.org == config_org)
    {
        *individual_config = new_config;
    } else {
        moonenv_config.profiles.push(new_config.clone());
    }

    MoonenvConfig::save_config(moonenv_config, default_path)?;

    Ok(())
}

/// Given a profile name, it will salve the file changing the default setting
pub fn set_config_name_as_default(name: String) -> Result<()> {
    let default_path = MoonenvConfig::get_default_config_file_path()?;
    let mut moonenv_config = MoonenvConfig::get_config(default_path.clone())?;

    moonenv_config.default = Some(name);

    MoonenvConfig::save_config(moonenv_config, default_path)?;

    Ok(())
}

/// Returns the default profile configuration
fn get_default_config() -> Result<IndividualConfig> {
    let moonenv_config = MoonenvConfig::get_config(MoonenvConfig::get_default_config_file_path()?)?;
    let default_config_name = moonenv_config
        .default
        .ok_or_else(|| anyhow::anyhow!("There is not default config set."))?;

    return get_config(default_config_name.borrow());
}

/// Returns the default profile name
fn get_default_org() -> Result<String> {
    let config = get_default_config()?;

    Ok(config.org)
}

/// Given a profile name, it will try to find the profile and return it
pub fn get_config(org: &str) -> Result<IndividualConfig> {
    let moonenv_config: MoonenvConfig =
        MoonenvConfig::get_config(MoonenvConfig::get_default_config_file_path()?)?;

    return moonenv_config
        .profiles
        .iter()
        .find(|config| config.org == *org)
        .ok_or_else(|| {
            anyhow::anyhow!(
                "No profile found. Ensure a default profile is correctly set in the configuration."
            )
        })
        .cloned();
}

/// If the org parameter is None, then it returns the default profile name;
/// Otherwise, it tries to return the given name
pub fn get_org(org: Option<String>) -> Result<String> {
    Ok(org
        .or(Some(get_default_org()?))
        .ok_or_else(|| anyhow::anyhow!("Org parameter is missing"))?)
}

/// Returns the API URL for a given profile name';
/// If the parameter is None, then it tries to return the default API URL;
pub fn get_url(org: &str) -> Result<String> {
    let config = get_config(org)?;

    Ok(config.url.ok_or_else(|| {
        anyhow::anyhow!("URL not configured in the profile. Please set the URL to proceed.")
    })?)
}

/// Returns the client ID for the given profile name;
/// If the parameter is None, then it tries to return the default client ID
pub fn get_client_id(org: &str) -> Result<String> {
    let config = get_config(org)?;

    Ok(config.client_id.ok_or_else(|| {
        anyhow::anyhow!(
            "Client Id is not configured in the profile. Please set the client ID to proceed"
        )
    })?)
}

#[cfg(test)]
mod test {
    use super::*;

    #[test]
    fn get_config_file_path_returns_the_config_file_path() {
        let result = MoonenvConfig::get_default_config_file_path();
        let mut root = home_dir().expect("It is expected a root directory path buf");
        root.push(".moonenv");
        root.push("config");

        assert_eq!(
            *result
                .expect("It is not expected an error")
                .to_str()
                .expect("It should return a path"),
            *root.to_str().expect("It should form the root path")
        )
    }

    #[test]
    fn get_config_file_returns_the_config() {
        let path = get_temp_path();

        let config = MoonenvConfig::get_config(path).expect("Failed to get config file");
        assert_eq!(config.default, Some("moonenv".to_string()));
        assert_eq!(config.profiles.len(), 1);
        assert_eq!(config.profiles[0].org, "moonenv");
    }

    #[test]
    fn save_config_writes_the_config() {
        let path = get_temp_path();
        let mut config = MoonenvConfig::default();

        config.default = Some("changed".to_string());

        let _ = MoonenvConfig::save_config(config, path.clone());
        let config: MoonenvConfig = confy::load_path(path).expect("Failed to load config file");

        assert_eq!(config.default, Some("changed".to_string()));
    }

    fn get_temp_path() -> PathBuf {
        let config_dir = tempfile::tempdir().expect("Failed to create tempfile");
        let mut path = config_dir.into_path();

        path.push(".moonenv");

        return path;
    }
}
