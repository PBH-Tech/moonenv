use anyhow::Result;
use confy::{self, ConfyError};
use serde::{Deserialize, Serialize};
use std::{env, path::PathBuf};

#[derive(Debug, Serialize, Deserialize)]
pub struct MoonenvConfig {
    pub default: Option<String>,

    pub profiles: Vec<IndividualConfig>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct IndividualConfig {
    pub url: Option<String>,

    pub name: String,
}

impl Default for IndividualConfig {
    fn default() -> Self {
        Self {
            name: "default".to_string(),
            url: Some("www.moonenv.app".to_string()),
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

fn get_config_path() -> PathBuf {
    let home = env::var("HOME").expect("HOME environment variable not set");
    let mut config_path: PathBuf = PathBuf::from(home);

    config_path.push(".moonenv");
    config_path.push("config");

    return config_path.clone();
}

fn get_config() -> Result<MoonenvConfig, ConfyError> {
    let config_path = get_config_path();

    return confy::load_path(config_path);
}

fn save_config(moonenv_config: MoonenvConfig) {
    let config_path = get_config_path();
    let _ = confy::store_path(config_path, moonenv_config);
}

pub fn change_config(new_config: IndividualConfig) -> Result<()> {
    let mut moonenv_config = get_config()?;
    let config_name = new_config.name.clone();

    if let Some(individual_config) = moonenv_config
        .profiles
        .iter_mut()
        .find(|config| config.name == config_name)
    {
        match &new_config.url {
            Some(url) => individual_config.url = Some(url.clone().to_string()),
            None => {}
        }
    } else {
        moonenv_config.profiles.push(new_config.clone());
    }

    save_config(moonenv_config);

    Ok(())
}

pub fn set_default(name: String) -> Result<()> {
    let mut moonenv_config = get_config()?;

    moonenv_config.default = Some(name);

    save_config(moonenv_config);

    Ok(())
}

fn get_default() -> Result<IndividualConfig> {
    let moonenv_config = get_config()?;

    return moonenv_config
        .profiles
        .iter()
        .find(|config| Some(config.name.to_string()) == moonenv_config.default)
        .ok_or_else(|| anyhow::anyhow!("No default profile found. Ensure a default profile is correctly set in the configuration.")).cloned();
}

pub fn get_default_url() -> Result<String> {
    let config = get_default()?;

    config.url.ok_or_else(|| {
        anyhow::anyhow!("URL not configured in the default profile. Please set the URL to proceed.")
    })
}
