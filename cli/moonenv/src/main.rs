use clap::{Parser, Subcommand, Args, ValueEnum};
/// Manages environment helping saving and pulling it
#[derive(Parser, Debug)]
pub struct App {
    #[clap(subcommand)]
    command: Command
}

#[derive(Subcommand, Debug)]
enum Command {
    /// Pulls the .env from the indicated repository
    Pull(RepoActionEnvArgs),

    /// Pushed the .env file located on the path where the command has been executed to the repository
    Push(RepoActionEnvArgs)
}

#[derive(Clone, ValueEnum, Debug)]
enum Environment {
    Dev,
    Qa,
    Prod,
}

#[derive(Args, Debug)]
struct RepoActionEnvArgs {
    /// Repository where to find the .env file
    repository: String,

    /// Environment where to find the .env file
    env: Environment
}



fn main() {
    let cli = App::parse();

    match cli.command {
        Command::Pull(value) => { pull_handler(value) }
        Command::Push(value) => { push_handler(value) }
    }

}

fn pull_handler(value: RepoActionEnvArgs) {
    println!("Pull: {:?}", value);
}

fn push_handler(value: RepoActionEnvArgs) {
    println!("Push: {:?}", value);
}
