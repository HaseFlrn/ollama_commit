# ollama_commit

## Overview

`ollama_commit` is a shell script that sets up an alias for a commit script in your shell configuration file. This allows you to easily run the commit script using the `commit` command.

## Prerequisites

Before you begin, ensure you have met the following requirements:

- You have a Unix-based operating system (Linux, macOS).
- You have a terminal application installed.
- You have write access to your shell configuration file (e.g., `.bashrc`, `.zshrc`).
- You have [Ollama](https://ollama.com/download) installed on your machine.

## Setup

To set up the `commit` alias, run the `setup` script. The script will check your shell configuration file and add the alias if it is not already present.

### Steps

1. Open a terminal.
2. Navigate to the directory containing the `setup` script.
3. Run the setup script:

   ```sh
   ./setup
   ```

4. If your shell configuration file is not found, you will need to pass it as an argument:

   ```sh
   ./setup /path/to/your/shell/config
   ```

## Usage

After running the setup script, source your shell configuration file to apply the changes:

```sh
source ~/.bashrc  # or ~/.zshrc, depending on your shell
```
