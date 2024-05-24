# GitHub Contributions TUI

This project is a terminal user interface (TUI) application built with the Bubble Tea framework. It fetches and displays GitHub contribution data for a specified user.

## Features

- Fetch GitHub contribution data using the GitHub GraphQL API.
- Display contribution data in a grid format.
- Simple and intuitive TUI for user interaction.

## Prerequisites

- Go (version 1.16 or later)
- A GitHub personal access token with the necessary permissions to access the GitHub GraphQL API.

## Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/jvanrhyn/github-contributions-tui.git
   cd github-contributions-tui
   ```

2. Create a `.env` file in the root directory of the project and add your GitHub token:

   ```sh
   echo "GITHUB_TOKEN=your_github_token" > .env
   ```

3. Install the dependencies:

   ```sh
   go mod tidy
   ```

## Usage

1. Run the application:

   ```sh
   go run main.go
   ```

2. Enter the GitHub username when prompted and press `Enter`.

3. The application will fetch and display the contribution data for the specified user.

4. To quit the application, press `Ctrl+C` or `esc`.

## Project Structure

- `main.go`: The main entry point of the application.
- `.env`: Environment file containing the GitHub token.
- `go.mod`: Go module file specifying the dependencies.
- `go.sum`: Go module file containing the cryptographic hashes of the dependencies.

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea): A powerful, elegant, and simple TUI framework for Go.
- [godotenv](https://github.com/joho/godotenv): A Go port of Ruby's dotenv library (loads environment variables from `.env`).

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [Charmbracelet](https://github.com/charmbracelet) for the Bubble Tea framework.
- [Joho](https://github.com/joho) for the godotenv library.
