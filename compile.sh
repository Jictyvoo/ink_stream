#!/bin/bash

# Define the input and output folders
INPUT_FOLDER="./extracted/Test"
OUTPUT_FOLDER="./output/Test"
KINDLE_VERSION=K11

# Run the kcc-c2e.py command
python3 /opt/kcc/kcc-c2e.py --profile $KINDLE_VERSION --manga-style -q --upscale --format auto --batchsplit 2 --output "$OUTPUT_FOLDER" $INPUT_FOLDER
