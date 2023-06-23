# DepViz

DepViz is a Go dependency visualization tool that generates a dependency tree for a Go project.

## Features

- Generates a JSON representation of the dependency tree for a Go project.
- Supports Go modules to accurately capture project dependencies.
- Uses native Go tools to ensure compatibility and ease of use.
- Outputs the dependency tree in a user-friendly JSON format.

## Installation

To install `depviz`, use the following command:

```bash
go get github.com/ramessesii2/depviz
```

or clone and build the project from source:

```bash
git clone https://www.github.com/ramessesi2/depviz
cd depviz
go build -o depviz ./cmd/cli/main.go
```

## Usage

To generate the dependency tree for a Go project, run the following command:

```bash
./depviz -r <repository-url> --ref <ref>
```

- `<repository-url>`: The FQDN of the Go project's GitHub repository.
- `<ref>`: The desired branch or tag of the repository.

The dependency tree will be generated and saved as a JSON file named `dependency_tree.json`.

## Example

Here's an example of how to use `depviz`:

```bash
./depviz -r  https://github.com/ramessesii2/depviz --ref main
```

This command generates the dependency tree for the `ramessesii2/depviz` repository, using the `main` branch.

The output file `dependency_tree.json` will contain the JSON representation of the dependency tree.
```json
[
  {
    "name": "github.com/ramessesii2/depviz",
    "version": "latest",
    "dependencies": [
      {
        "name": "github.com/inconshreveable/mousetrap",
        "version": "v1.1.0",
        "dependencies": []
      },
      {
        "name": "github.com/spf13/cobra",
        "version": "v1.7.0",
        "dependencies": [
          {
            "name": "github.com/cpuguy83/go-md2man/v2",
            "version": "v2.0.2",
            "dependencies": []
          },
          {
            "name": "github.com/inconshreveable/mousetrap",
            "version": "v1.1.0",
            "dependencies": []
          },
          {
            "name": "github.com/spf13/pflag",
            "version": "v1.0.5",
            "dependencies": []
          },
          {
            "name": "gopkg.in/yaml.v3",
            "version": "v3.0.1",
            "dependencies": []
          }
        ]
      },
      {
        "name": "github.com/spf13/pflag",
        "version": "v1.0.5",
        "dependencies": []
      }
    ]
  }
]
```

## Contributing

Contributions are welcome! If you encounter any issues or have suggestions for improvements, please open an issue or submit a pull request on the GitHub repository.

## License

This project is licensed under the [GNU GENERAL PUBLIC LICENSE](LICENSE).

---
