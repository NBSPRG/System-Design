# LLD Repository Makefile
# Provides convenient commands to build and run projects

.PHONY: help setup clean build-all demo-all

# Default target
help:
	@echo "LLD Repository - Available Commands:"
	@echo ""
	@echo "Setup & Management:"
	@echo "  make setup         - Install dependencies for all projects"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make build-all     - Build all projects"
	@echo ""
	@echo "Project Demos:"
	@echo "  make demo-bloom    - Run Bloom Filter demo"
	@echo "  make demo-circuit  - Run Circuit Breaker demo"
	@echo "  make demo-hashmap  - Run HashMap demo"
	@echo "  make demo-tictac   - Run Tic Tac Toe game"
	@echo "  make demo-library  - Run Library Management demo"
	@echo ""
	@echo "Quick Start:"
	@echo "  make setup && make demo-bloom"

# Setup dependencies
setup:
	@echo "Setting up LLD repository..."
	@echo "Installing Go dependencies..."
	cd design-patterns/bloom-filter && go mod tidy
	cd design-patterns/circuit-breaker && go mod tidy
	@echo "Setup complete!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	find . -name "*.exe" -delete
	find . -name "a.out" -delete
	find . -name "__pycache__" -type d -exec rm -rf {} +
	@echo "Clean complete!"

# Build all projects
build-all:
	@echo "Building all projects..."
	# Go projects
	cd design-patterns/bloom-filter && go build
	cd design-patterns/circuit-breaker && make build
	# C++ projects
	cd algorithms/hashmap && g++ -o hashmap hashmap.cpp
	cd games/tic-tac-toe && g++ -o game main.cpp Board/Board.cpp Game/Game.cpp Player/Player.cpp
	cd games/snake-ladder && g++ -o snake-ladder snake-ladder.cpp
	cd systems/split-wise && g++ -o split-wise split-wise.cpp
	@echo "Build complete!"

# Demo commands
demo-bloom:
	@echo "Running Bloom Filter demo..."
	cd design-patterns/bloom-filter && go run .

demo-circuit:
	@echo "Running Circuit Breaker demo..."
	cd design-patterns/circuit-breaker && make demo

demo-hashmap:
	@echo "Running HashMap demo..."
	cd algorithms/hashmap && g++ -o hashmap hashmap.cpp && ./hashmap

demo-tictac:
	@echo "Running Tic Tac Toe game..."
	cd games/tic-tac-toe && g++ -o game main.cpp Board/Board.cpp Game/Game.cpp Player/Player.cpp && ./game

demo-library:
	@echo "Running Library Management demo..."
	cd systems/library-management && python Demo.py

demo-scheduling:
	@echo "Running Scheduling demo..."
	cd algorithms/scheduling && python load_balancer.py

# Run all demos sequentially
demo-all: demo-bloom demo-hashmap demo-library

# Development helpers
list-projects:
	@echo "Available projects:"
	@find . -name "README.md" -not -path "./README.md" -exec dirname {} \; | sort

check-structure:
	@echo "Repository structure:"
	@tree -d -L 3 || find . -maxdepth 3 -type d | sort

# Project-specific builds
build-go:
	cd design-patterns/bloom-filter && go build
	cd design-patterns/circuit-breaker && make build

build-cpp:
	cd algorithms/hashmap && g++ -o hashmap hashmap.cpp
	cd games/tic-tac-toe && g++ -o game main.cpp Board/Board.cpp Game/Game.cpp Player/Player.cpp
	cd games/snake-ladder && g++ -o snake-ladder snake-ladder.cpp
	cd systems/split-wise && g++ -o split-wise split-wise.cpp
