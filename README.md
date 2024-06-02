<div>
  <picture>
    <img alt="Warp" width="230" src="https://github.com/PBH-Tech/moonenv/assets/moonenv.png">
  </picture>
</div>
<br />
<div>
<a href="https://buymeacoffee.com/moonenv" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;" ></a>
<div/>

# Moonenv

Moonenv is a Git-like command-line tool designed to help developers manage their environment variable files with ease. Inspired by the simplicity and power of Git, Moonenv provides a robust way to handle environment settings across different development stages, ensuring consistency and reducing errors.

## Philosophy

Moonenv is built on the belief that developers should have full control and ownership of their environment configurations, reflecting our core philosophies:

- **Ownership of Environment Files:** Your environment settings are as crucial as your codebase. With Moonenv, the env files belong entirely to you. This is especially relevant in infrastructure-as-code (IaC) setups where you can spin up your own instances and manage them directly. Moonenv empowers you to maintain and control your environment variables independently.

- **No Cost Necessary:** We believe that essential tools for software development, like environment variable management, should be accessible to everyone. Moonenv is free to use, ensuring that you don't have to pay to maintain your environment settings securely and efficiently.

- **Open Source Commitment:** Moonenv is an open-source project. We are committed to maintaining transparency, promoting community contributions, and ensuring that our tool can be trusted and improved upon by the community. We encourage you to dive into the code, provide feedback, and contribute to the project.

This philosophy drives every decision we make for Moonenv and aims to support developers in creating secure, efficient, and scalable applications without the overhead of managing and syncing environment variables traditionally.

## Installation

Binaries way is not available yet. Follow the issue [Binary distribution](https://github.com/PBH-Tech/moonenv/issues/13) for more information.
 
### From source

With Rust's package manager cargo, you can install Moonenv via:

```bash
cargo install moonenv
```

## Getting Started

### Configuration

When you first use Moonenv in a project, it creates a `/.moonenv/config` file in the root directory. This configuration file is where all your settings live, including the ability to define multiple profiles and set a default profile. For ease of management and to avoid manual editing errors, we strongly recommend using the `moonenv config` command to manipulate these profiles.

### Commands

Moonenv is designed to be intuitive for those familiar with Git-like command structures. Every command in Moonenv includes a `--help` option, which you can use to understand what the options for each command are. Below are some of the primary commands available:

- **config**: Changes the application's configuration settings. Use this to upsert profiles and set default behaviors.
- **pull**: Pulls the current .env file from the specified repository. This is useful for syncing the latest environment settings that have been pushed to your repository.
- **push**: Pushes the local .env file from the directory where the command is executed to the repository. This ensures that your remote settings are up-to-date.

## Getting Help

To get more information on any command, you can always run:

```bash
moonenv <command> --help
```

This will display detailed usage instructions for the specified command, helping you understand all available options and how to use them effectively.

## Contributing

Contributions are welcome! Please feel free to fork the repository, make your changes, and submit a pull request.