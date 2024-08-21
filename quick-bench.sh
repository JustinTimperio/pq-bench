#!/bin/bash

# Go
echo "=============================="
echo "Compiling and running Go program..."
cd ./go
go build .
time ./bench
cd ..
echo ""

## GPQ
echo "=============================="
echo "Compiling and running GPQ program..."
cd ./gpq
go build .
time ./bench
cd ..
echo ""

# Zig
echo "=============================="
echo "Compiling and running Zig program..."
cd ./zig
zig build -Doptimize=ReleaseFast
time ./zig-out/bin/zig
cd ..
echo ""

# Rust
echo "=============================="
echo "Compiling and running Rust program..."
cd ./rust
cargo build --release
time ./target/release/bench
cd ..
echo ""

# C++
echo "=============================="
echo "Compiling and running C++ program..."
cd ./c++
g++ -O3 bench.cpp
time ./a.out
cd ..
