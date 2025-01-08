use anyhow::{anyhow, Ok, Result};
use confy::{self};
use dirs::home_dir;
use serde::{Deserialize, Serialize};
use std::{borrow::Borrow, path::PathBuf};

#[derive(Debug, Serialize, Deserialize)]
pub struct MoonenvConfig {
    pub default: Option<String>,

    pub profiles: Vec<IndividualConfig>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct IndividualConfig {
    pub url: Option<String>,

    pub org: String,

    pub device_code: Option<String>,

    pub access_token: Option<String>,

    pub refresh_token: Option<String>,

    pub client_id: Option<String>,
}

impl Default for IndividualConfig {
    fn default() -> Self {
        Self {
            org: "moonenv".to_string(),
            url: Some("www.moonenv.app".to_string()),
            access_token: None,
            device_code: None,
            refresh_token: None,
            client_id: None,
        }
    }
}

impl ::std::default::Default for MoonenvConfig {
    fn default() -> Self {
        Self {
            default: Some("default".into()),
            profiles: vec![IndividualConfig::default()],
        }
    }
}

fn get_config_file_path() -> Result<PathBuf> {
    let home = home_dir().ok_or_else(|| anyhow::anyhow!("HOME environment variable not set"))?;
    let mut config_path: PathBuf = PathBuf::from(home);

    config_path.push(".moonenv");
    config_path.push("config");

    Ok(config_path.clone())
}

fn get_config_file() -> Result<MoonenvConfig> {
    let config_path = get_config_file_path()?;

    return confy::load_path(config_path)
        .map_err(|e| anyhow!("Failed to load configuration: {}", e));
}

fn save_config(moonenv_config: MoonenvConfig) -> Result<()> {
    let config_path: PathBuf = get_config_file_path()?;
    let _ = confy::store_path(config_path, moonenv_config);

    Ok(())
}

pub fn change_config(new_config: IndividualConfig) -> Result<()> {
    let mut moonenv_config = get_config_file()?;
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

    let _ = save_config(moonenv_config)?;

    Ok(())
}

pub fn set_config_name_as_default(name: String) -> Result<()> {
    let mut moonenv_config = get_config_file()?;

    moonenv_config.default = Some(name);

    let _ = save_config(moonenv_config)?;

    Ok(())
}

fn get_default_config() -> Result<IndividualConfig> {
    let moonenv_config = get_config_file()?;
    let default_config_name = moonenv_config
        .default
        .ok_or_else(|| anyhow::anyhow!("There is not default config set."))?;

    return get_config(default_config_name.borrow());
}

fn get_default_org() -> Result<String> {
    let config = get_default_config()?;

    Ok(config.org)
}

pub fn get_config(org: &String) -> Result<IndividualConfig> {
    let moonenv_config = get_config_file()?;

    return moonenv_config
        .profiles
        .iter()
        .find(|config| config.org.to_string() == *org)
        .ok_or_else(|| {
            anyhow::anyhow!(
                "No profile found. Ensure a default profile is correctly set in the configuration."
            )
        })
        .cloned();
}

pub fn get_org(org: Option<String>) -> Result<String> {
    Ok(org
        .or(Some(get_default_org()?))
        .ok_or_else(|| anyhow::anyhow!("Org parameter is missing"))?)
}

pub fn get_url(org: &String) -> Result<String> {
    let config = get_config(org)?;

    Ok(config.url.ok_or_else(|| {
        anyhow::anyhow!("URL not configured in the profile. Please set the URL to proceed.")
    })?)
}

pub fn get_client_id(org: &String) -> Result<String> {
    let config = get_config(org)?;

    Ok(config.client_id.ok_or_else(|| {
        anyhow::anyhow!(
            "Client Id is not configured in the profile. Please set the client ID to proceed"
        )
    })?)
}
