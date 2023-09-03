#!/bin/bash

# Exit immediately if a command exits with a non-zero status and print command traces before executing the command
set -e
set -x

# Initialize variables
tenantName="testtenant3"
target="testtarget3"
targetGroup="testtargetgroup3"
targetNode="edgepc2"
outputDir="./outputCredentials"

# Setup and teardown functions
setUp() {
    rm -rf $outputDir
    mkdir -p $outputDir
}

tearDown() {
    echo "Cleaning up..."
    rm -rf $outputDir
}

# Test functions

test_create_tenant() {
    echo "Test: Creating tenant $tenantName"
    ./kufast create tenant $tenantName --cpu="500m" --memory="1Gi" --pods="1" --output=$outputDir
}

test_create_nginx_pod() {
    echo "Test: Creating Nginx Pod"
    ./kufast create pod mypod nginx
    echo "yes" | ./kufast delete pod mypod
}

test_delete_tenant() {
    echo "Test: Deleting tenant $tenantName"
    echo "yes" | ./kufast delete tenant $tenantName
}

test_create_and_delete_secret() {
    echo "Test: Creating and Deleting Secret"
    echo "password" | ./kufast create secret credentials
    echo "yes" |  ./kufast delete secret credentials
}

test_create_and_delete_target_group() {
    echo "Test: Creating and Deleting Target Group"
    ./kufast create target-group $targetGroup $targetNode
    echo "yes" | ./kufast delete target-group $targetGroup $targetNode
}

# Main execution

setUp

# Execute tests and store results
testResults=()

echo "Executing tests..."

test_create_tenant && testResults+=("test_create_tenant: Passed") || testResults+=("test_create_tenant: Failed")
test_create_nginx_pod && testResults+=("test_create_nginx_pod: Passed") || testResults+=("test_create_nginx_pod: Failed")
test_delete_tenant && testResults+=("test_delete_tenant: Passed") || testResults+=("test_delete_tenant: Failed")
test_create_and_delete_secret && testResults+=("test_create_and_delete_secret: Passed") || testResults+=("test_create_and_delete_secret: Failed")
test_create_and_delete_target_group && testResults+=("test_create_and_delete_target_group: Passed") || testResults+=("test_create_and_delete_target_group: Failed")

# Summary
echo "Test summary:"
for result in "${testResults[@]}"; do
    echo "  - $result"
done

tearDown
