#!/bin/sh

# run the server using
# go run main.go server run -d --memory_percentage=0.9 --estimate_with_internal_memory=false


go run main.go client compare --run_original=true  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_first_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_896x896_1.0
go run main.go client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_second_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_896x896_1.0
go run main.go client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_third_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_896x896_1.0
go run main.go client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_forth_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_896x896_1.0

go run main.go client compare --run_original=true  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_first_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_681x681_1.0
go run main.go client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_second_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_681x681_1.0
go run main.go client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_third_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_681x681_1.0
go run main.go client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_forth_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_681x681_1.0

go run main.go client compare --run_original=true  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_first_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_672x672_1.0
go run main.go client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_second_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_672x672_1.0
go run main.go client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_third_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_672x672_1.0
go run main.go client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_forth_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_672x672_1.0

go run main.go client compare --run_original=true  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_first_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_908x908_1.0
go run main.go client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_second_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_908x908_1.0
go run main.go client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_third_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_908x908_1.0
go run main.go client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_forth_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_908x908_1.0

go run main.go client compare --run_original=true  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_first_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_448x448_1.0
go run main.go client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_second_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_448x448_1.0
go run main.go client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_third_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_448x448_1.0
go run main.go client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_forth_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_448x448_1.0

go run main.go client compare --run_original=true  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_first_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_227x227_1.0
go run main.go client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_second_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_227x227_1.0
go run main.go client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_third_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_227x227_1.0
go run main.go client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_forth_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_227x227_1.0

go run main.go client compare --run_original=true  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_first_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_224x224_1.0
go run main.go client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_second_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_224x224_1.0
go run main.go client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_third_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_224x224_1.0
go run main.go client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_forth_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_224x224_1.0
