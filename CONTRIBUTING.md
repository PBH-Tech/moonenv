# Contributing to Moonenv

> Thank you very much for considering contributing to Moonenv. Feel free to share your ideas and suggest changes (always respecting the Code of conduct). 😄

Moonenv is an open-source CLI tool that helps developers manage, version, and share their environment variables. You can contribute in many ways: helping us with documentation, submitting bug reports, requesting new features or writing new Code that can be merged into and used in future versions.

## Project structure

Moonenv is split into CLI in Rust and the Server in Go.

> ⚠ Be aware that we are not experts in any of these langues; This project is also our process of learning new programming languages. Feel free to help us to improve it or alert us if you find any vulnerability.

### Server

For the Server deployment, we use [CDK](https://docs.aws.amazon.com/cdk/v2/guide/home.html) as our Infra as Code (IaC), and it has the following stacks:

- MoonenvRoute53Stack: Responsible for keeping the Route 53 and certificates resources used for all the other stacks;
- MoonenvS3Stack: Contains the Bucket S3 resource where the environment variable is saved;
- MoonenvDynamoDb: Stack that creates the dynamo DB tables;
- MoonenvCognitoStack: Cognito is the Auth solution for Moonenv; this stack sets this service;
- MoonenvLambdaStack: Contains all the necessary lambdas for the server;
- MoonenvApiGatewayStack: The Rest API Gateway service that the CLI communicates with;

We decided on as many serverless solutions as possible because we believe that environment resources are not often requested, so serverless is the cheaper solution.

Before deploying it, create a `.env` file inside the server folder and get the example in the `server/.env.example` file.

`cdk deploy --all` is the command for the deployment, but we strongly recommend that you to read more about CDK if you have no experience.

## Reporting bugs, requesting new features and improving docs

Creating an issue or requesting a pull request are the best ways you can contribute. Please fill out the templates and provide as much information as possible to make it clear for us to understand. And one more time, thank you very much for helping us to build this project. ♥  