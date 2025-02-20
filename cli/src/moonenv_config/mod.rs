///! Moonenv configuration module
/// This module is responsible for handling the configuration for each profile.
/// It contains the MoonenvConfig struct, which is responsible for storing the default profile and all profiles that are available.
///
use anyhow::Result;
use file_handler::FileHandler;
use serde::{Deserialize, Serialize};
use std::time::Duration;

mod file_handler;

/// The configuration struct, responsible to informs the default profile and all profiles that are available
#[derive(Debug, Serialize, Deserialize, Clone)]
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

impl MoonenvConfig {
    pub fn new() -> Self {
        let moonenv_config = Self::get_file(Self::get_default_file_path().unwrap()).unwrap();

        Self {
            default: moonenv_config.default,
            profiles: moonenv_config.profiles,
        }
    }

    /// Change an individual configuration;
    /// If it already exists, so it overwrites;
    /// Otherwise, it pushes a new profile in the profile list;
    /// And finally, it saves the config in the file
    pub fn change_config(&mut self, new_config: IndividualConfig) -> Result<()> {
        if let Some(individual_config) = self
            .profiles
            .iter_mut()
            .find(|config| config.org == new_config.org)
        {
            *individual_config = new_config;
        } else {
            self.profiles.push(new_config);
        }

        Self::save_file(self.clone(), Self::get_default_file_path()?)?;

        Ok(())
    }

    /// Given a profile name, it will salve the file changing the default setting
    pub fn set_config_name_as_default(&mut self, name: String) -> Result<()> {
        let default_path = Self::get_default_file_path()?;

        self.default = Some(name);

        Self::save_file(self.clone(), default_path)?;

        Ok(())
    }

    /// Given a profile name, it will try to find the profile and return it
    pub fn get_config(&mut self, org: &str) -> Result<IndividualConfig> {
        return self
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
    pub fn get_org(&mut self, org: Option<String>) -> Result<String> {
        Ok(org
            .or_else(|| self.default.clone())
            .ok_or_else(|| anyhow::anyhow!("Org parameter is missing"))?)
    }

    /// Returns the API URL for a given profile name';
    /// If the parameter is None, then it tries to return the default API URL;
    pub fn get_url(&mut self, org: &str) -> Result<String> {
        let config = self.get_config(org)?;

        Ok(config.url.ok_or_else(|| {
            anyhow::anyhow!("URL not configured in the profile. Please set the URL to proceed.")
        })?)
    }

    /// Returns the client ID for the given profile name;
    /// If the parameter is None, then it tries to return the default client ID
    pub fn get_client_id(&mut self, org: &str) -> Result<String> {
        let config = self.get_config(org)?;

        Ok(config.client_id.ok_or_else(|| {
            anyhow::anyhow!(
                "Client Id is not configured in the profile. Please set the client ID to proceed"
            )
        })?)
    }
}

impl FileHandler<MoonenvConfig> for MoonenvConfig {}
