#!/bin/sh

# run the server using
# ./main server run -d --memory_percentage=0.9 --estimate_with_internal_memory=false

killall uprd 2>/dev/null
killall main 2>/dev/null

sleep 10


go build main.go

echo "##########################################################"
echo "##### VGG 896"
echo "##########################################################"

./main server run --memory_percentage=0.9 --estimate_with_internal_memory=false &

export UPRD_PID=$!

sleep 10

./main client compare --run_original=true  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_first_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_896x896_1.0
./main client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_second_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_896x896_1.0
./main client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_third_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_896x896_1.0
./main client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_forth_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_896x896_1.0

kill -9 $UPRD_PID
killall uprd 2>/dev/null
killall main 2>/dev/null

sleep 10

##########################################################
##########################################################

echo "##########################################################"
echo "##### ALEXNET 681"
echo "##########################################################"

./main server run --memory_percentage=0.9 --estimate_with_internal_memory=false &

export UPRD_PID=$!

sleep 10

./main client compare --run_original=true  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_first_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_681x681_1.0
./main client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_second_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_681x681_1.0
./main client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_third_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_681x681_1.0
./main client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_forth_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_681x681_1.0

kill -9 $UPRD_PID
killall uprd 2>/dev/null
killall main 2>/dev/null

sleep 10

##########################################################
##########################################################

echo "##########################################################"
echo "##### VGG 672"
echo "##########################################################"

./main server run --memory_percentage=0.9 --estimate_with_internal_memory=false &

export UPRD_PID=$!

sleep 10

./main client compare --run_original=true  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_first_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_672x672_1.0
./main client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_second_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_672x672_1.0
./main client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_third_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_672x672_1.0
./main client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_forth_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_672x672_1.0

kill -9 $UPRD_PID

killall uprd 2>/dev/null
killall main 2>/dev/null

sleep 10

##########################################################
##########################################################

echo "##########################################################"
echo "##### ALEXNET 908"
echo "##########################################################"

./main server run --memory_percentage=0.9 --estimate_with_internal_memory=false &

export UPRD_PID=$!

sleep 10

./main client compare --run_original=true  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_first_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_908x908_1.0
./main client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_second_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_908x908_1.0
./main client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_third_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_908x908_1.0
./main client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_forth_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_908x908_1.0


kill -9 $UPRD_PID

killall uprd 2>/dev/null
killall main 2>/dev/null

sleep 10

##########################################################
##########################################################

echo "##########################################################"
echo "##### VGG 448"
echo "##########################################################"

./main server run --memory_percentage=0.9 --estimate_with_internal_memory=false &

export UPRD_PID=$!

sleep 10

./main client compare --run_original=true  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_first_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_448x448_1.0
./main client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_second_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_448x448_1.0
./main client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_third_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_448x448_1.0
./main client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_forth_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_448x448_1.0

kill -9 $UPRD_PID

killall uprd 2>/dev/null
killall main 2>/dev/null

sleep 10

##########################################################
##########################################################

echo "##########################################################"
echo "##### ALEXNET 454"
echo "##########################################################"

./main server run --memory_percentage=0.9 --estimate_with_internal_memory=false &

export UPRD_PID=$!

sleep 10

./main client compare --run_original=true  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_first_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_454x454_1.0
./main client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_second_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_454x454_1.0
./main client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_third_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_454x454_1.0
./main client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_forth_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_454x454_1.0

kill -9 $UPRD_PID

killall uprd 2>/dev/null
killall main 2>/dev/null

sleep 10

##########################################################
##########################################################

echo "##########################################################"
echo "##### ALEXNET 227"
echo "##########################################################"

./main server run --memory_percentage=0.9 --estimate_with_internal_memory=false &

export UPRD_PID=$!

sleep 10

./main client compare --run_original=true  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_first_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_227x227_1.0
./main client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_second_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_227x227_1.0
./main client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_third_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_227x227_1.0
./main client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_forth_iteration --experiment_description="large model comparison" --large_models=true --models=large_alexnet_227x227_1.0

kill -9 $UPRD_PID

killall uprd 2>/dev/null
killall main 2>/dev/null

sleep 10

##########################################################
##########################################################

echo "##########################################################"
echo "##### VGG 224"
echo "##########################################################"

./main server run --memory_percentage=0.9 --estimate_with_internal_memory=false &

export UPRD_PID=$!

sleep 10

./main client compare --run_original=true  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_first_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_224x224_1.0
./main client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_second_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_224x224_1.0
./main client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_third_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_224x224_1.0
./main client compare --run_original=false  --monitor_memory=false --concurrent=1  --profile_output=large_model_compare_forth_iteration --experiment_description="large model comparison" --large_models=true --models=large_vgg16_224x224_1.0

kill -9 $UPRD_PID

killall uprd 2>/dev/null
killall main 2>/dev/null

sleep 10

##########################################################
##########################################################

echo "done running large models"
