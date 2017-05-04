#!/usr/bin/bash

set -e 

source /cluster_share/OHSU/shared/.venv/bin/activate

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m'
function demoRun() {
    echo -e "${YELLOW}$@ --print${NC}"
    echo
    eval "$@ --print | jq ." 
    echo
    read -n 1 -s -p "Press any key to continue..." key
    echo
    eval "$@ --wait"
    echo
}

function demoCmd() {
    echo -e "${YELLOW}$@${NC}"
    echo
    eval "$@"
    echo 
    read -n 1 -s -p "Press any key to continue..." key
    echo
}


DEMO_DIR=/cluster_share/home/strucka/funnel_demo
mkdir -p $DEMO_DIR $DEMO_DIR/routed_file $DEMO_DIR/fetch_file $DEMO_DIR/push_file


echo -e "${GREEN}
##-----------------------------------------##
## ROUTED FILE TEST
##-----------------------------------------##
${NC}"
TEST_FILE=${DEMO_DIR}/routed_file/test_local_input
TEST_OUTFILE=${DEMO_DIR}/routed_file/test_local_output

# Stage run
touch $TEST_FILE
echo 'LOCAL FILE' > $TEST_FILE
CCCID=$(ccc_client dts post -f $TEST_FILE -s ohsu -u strucka | cut -f 2)

demoRun funnel run "'md5sum \$INFILE > \$OUTFILE'" \
--server http://application-0-1:18000 \
--container ubuntu \
--in INFILE=ccc://$CCCID \
--out OUTFILE=ccc://$TEST_OUTFILE \
--tag strategy=routed_file

# Check output
echo -e "${NC}Task Working Directory:${NC}"
demoCmd ls -lah $DEMO_DIR/routed_file
echo -e "${NC}Output DTS Recod:${NC}"
demoCmd "ccc_client dts get $TEST_OUTFILE | jq ."


echo -e "${GREEN}
##-----------------------------------------##
## FETCH FILE TEST
##-----------------------------------------##
${NC}"
TEST_FILE=${DEMO_DIR}/fetch_file/test_remote_input
TEST_OUTFILE=${DEMO_DIR}/fetch_file/test_local_output

# Stage run
touch $TEST_FILE
echo 'REMOTE FILE' > $TEST_FILE
scp $TEST_FILE central-gateway:$TEST_FILE
CCCID=$(ccc_client dts post -f $TEST_FILE -s ohsu -u strucka | cut -f 2)
rm $TEST_FILE

demoRun funnel run "'md5sum \$INFILE > \$OUTFILE'" \
--server http://application-0-1:18000 \
--container ubuntu \
--in INFILE=ccc://$CCCID \
--out OUTFILE=ccc://$TEST_OUTFILE \
--tag strategy=fetch_file

# Input should have 2 locations now
demoCmd ccc_client dts get $(cat $TEST_FILE_CCCID)

# Check output
echo -e "${NC}Task Working Directory:${NC}"
demoCmd ls -lah $DEMO_DIR/fetch_file
echo -e "${NC}Output DTS Recod:${NC}"
demoCmd "ccc_client dts get $TEST_OUTFILE | jq ."


echo -e "${GREEN}
##-----------------------------------------##
## PUSH FILE TEST
##-----------------------------------------##
${NC}"
TEST_FILE=${DEMO_DIR}/push_file/test_local_input
TEST_OUTFILE=${DEMO_DIR}/push_file/test_pushed_output

# Stage run
touch $TEST_FILE
echo 'LOCAL FILE' > $TEST_FILE
CCCID=$(ccc_client dts post -f $TEST_FILE -s ohsu -u strucka | cut -f 2)

demoRun funnel run "'md5sum \$INFILE > \$OUTFILE'" \
--server http://application-0-1:18000 \
--container ubuntu \
--in INFILE=ccc://$CCCID \
--out OUTFILE=ccc://$TEST_OUTFILE \
--tag strategy=push_file

# Check output
echo -e "${NC}Task Working Directory:${NC}"
demoCmd ls -lah $DEMO_DIR/push_file
echo -e "${NC}Output DTS Recod:${NC}"
demoCmd "ccc_client dts get $TEST_OUTFILE | jq ."

