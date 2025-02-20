//! File Handler
//!
//! This module is responsible for saving and reading the config in the ~/.moonenv/config file.
//! It contains functions that validates and gives the configuration file.

use anyhow::{anyhow, Ok, Result};
use confy::{self};
use dirs::home_dir;
use serde::{de::DeserializeOwned, Serialize};
use std::path::PathBuf;

/// The FileHandler
pub trait FileHandler<T>
where
    T: DeserializeOwned + Default + Serialize,
{
    /// Returns the file path for the configuration file
    fn get_default_file_path() -> Result<PathBuf> {
        let mut home = home_dir().ok_or_else(|| anyhow!("HOME environment variable not set"))?;

        home.push(".moonenv");
        home.push("config");

        Ok(home)
    }

    /// Saves the configuration
    fn save_file(moonenv_config: T, path: PathBuf) -> Result<()> {
        let _ = confy::store_path(path, moonenv_config);

        Ok(())
    }

    /// Returns the configuration
    fn get_file(config_path: PathBuf) -> Result<T> {
        confy::load_path(config_path).map_err(|e| anyhow!("Failed to load configuration: {}", e))
    }
}
