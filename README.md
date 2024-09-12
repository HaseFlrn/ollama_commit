# Ollama Commit

Ollama Commit is a tool designed to help developers generate concise and descriptive commit messages based on the changes in their Git repository. By leveraging a local AI model, it ensures that commit messages follow a specific format and are both informative and standardized.

## Features

- **Automated Commit Message Generation**: Generates commit messages based on the diff of staged changes.
- **Configurable AI Model**: Uses a local AI model to generate commit messages.
- **Prerequisite Checks**: Ensures that necessary tools like Git and Ollama are installed and running.
- **Interactive Confirmation**: Prompts the user for confirmation before committing changes.

## Installation

> [!NOTE]
> Currently only the self-installation via go-cli and manually setting the PATH is available. This will change in the future. ([see here](./TODO.md))

1. **Clone the Repository**:

   ```sh
   git clone https://github.com/HaseFlrn/ollama_commit.git
   cd ollama_commit
   ```

2. **Install Dependencies**:
   Ensure you have Go installed. Then, run:

   ```sh
   go mod tidy
   ```

3. **Build the Project**:

   ```sh
   go build -o ollama_commit
   ```

## Usage

1. **Start Ollama**:
   Ensure that the Ollama service is running on `http://localhost:11434/`.

2. **Run the Tool**:
   Navigate to your Git repository and run:

   ```sh
   /path/to/ollama_commit
   ```

3. **Follow the Prompts**:
   - The tool will check if you are inside a Git repository.
   - It will then generate a commit message based on the staged changes.
   - You will be prompted to confirm the commit message before it is committed to the repository.

## Example

```sh
cd /path/to/your/git/repo
/path/to/ollama_commit
```

## Configuration

Currently, the AI model used for generating commit messages is hardcoded to llama3. Future versions may include options to configure the model and other parameters.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## License

This project is licensed under the MIT License. See the LICENSE file for details.
