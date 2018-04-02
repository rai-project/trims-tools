#!/bin/sh

# run server using go run main.go server run -d --memory_percentage=0.8 --monitor_memory=false

go run main.go client compare --monitor_memory=false --concurrent=1 --run_original=true  --combine=false --profile_output=first_comparison --experiment_description="model comparison"  --profile_output_overwrite=true
go run main.go client compare --monitor_memory=false --concurrent=1 --run_original=false --combine=false --profile_output=second_comparison --experiment_description="model comparison"  --profile_output_overwrite=true
go run main.go client compare --monitor_memory=false --concurrent=1 --run_original=false --combine=false --profile_output=third_comparison --experiment_description="model comparison"  --profile_output_overwrite=true
go run main.go client compare --monitor_memory=false --concurrent=1 --run_original=false --combine=false --profile_output=forth_comparison --experiment_description="model comparison"  --profile_output_overwrite=true
