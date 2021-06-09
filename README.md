# Package Typospotter

This project aims to provide a basis for identifying packages potentially [typosquatting](https://snyk.io/blog/typosquatting-attacks/) within package repository at the point of addition to the package repository.

**N.B** The current implementation is quite basic and should not be considered production ready.

## Usage

```
docker build -t package-typospotter:latest .
docker run -it package-typospotter:latest

# Alternatively

go build -o package-typospotter cmd/squatter-spotter/main.go
./package-typospotter
```
