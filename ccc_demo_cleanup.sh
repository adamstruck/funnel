#!/usr/bin/bash

set -e 

source /cluster_share/OHSU/shared/.venv/bin/activate

DEMO_DIR=/cluster_share/home/strucka/funnel_demo

##-----------------------------------------##
## ROUTED FILE TEST
##-----------------------------------------##
TEST_FILE=${DEMO_DIR}/routed_file/test_local_input
TEST_OUTFILE=${DEMO_DIR}/routed_file/test_local_output

# Cleanup previous runs
if [ -e $TEST_FILE ]; then
    ID=$(ccc_client dts infer-cccId $TEST_FILE | cut -f 2)
    ccc_client dts delete $ID
    rm $TEST_FILE
fi
if [ -e $TEST_OUTFILE ]; then
    ccc_client dts delete $TEST_OUTFILE
    rm $TEST_OUTFILE
fi
##-----------------------------------------##
## PUSH FILE TEST
##-----------------------------------------##
TEST_FILE=${DEMO_DIR}/push_file/test_local_input
TEST_OUTFILE=${DEMO_DIR}/push_file/test_pushed_output

# Cleanup previous runs
if [ -e $TEST_FILE ]; then
    ID=$(ccc_client dts infer-cccId $TEST_FILE | cut -f 2)
    ccc_client dts delete $ID
    rm $TEST_FILE
fi
if [ -e $TEST_OUTFILE ]; then
    ccc_client dts delete $TEST_OUTFILE
    rm $TEST_OUTFILE
fi

##-----------------------------------------##
## FETCH FILE TEST
##-----------------------------------------##
TEST_FILE=${DEMO_DIR}/fetch_file/test_remote_input
TEST_OUTFILE=${DEMO_DIR}/fetch_file/test_local_output

# Cleanup previous runs
if [ -e $TEST_FILE ]; then
    ID=$(ccc_client dts infer-cccId $TEST_FILE | cut -f 2)
    ccc_client dts delete $ID
    rm $TEST_FILE
fi
if [ -e $TEST_OUTFILE ]; then
    ccc_client dts delete $TEST_OUTFILE
    rm $TEST_OUTFILE
fi
