#!/bin/sh

# run server using go run main.go server run -d --memory_percentage=0.8 --persist_cpu_only=true --monitor_memory=false

go run main.go client compare --monitor_memory=false --concurrent=1  --profile_output=first_cpu_only_comparison --experiment_description="model comparison with only cpu persistence"  --profile_output_overwrite=true
go run main.go client compare --monitor_memory=false --concurrent=1  --profile_output=second_cpu_only_comparison --experiment_description="model comparison with only cpu persistence"  --profile_output_overwrite=true
#go run main.go client compare --monitor_memory=false --concurrent=1  --profile_output=third_cpu_only_comparison --experiment_description="model comparison with only cpu persistence"  --profile_output_overwrite=true
#go run main.go client compare --monitor_memory=false --concurrent=1  --profile_output=forth_cpu_only_comparison --experiment_description="model comparison with only cpu persistence"  --profile_output_overwrite=true
