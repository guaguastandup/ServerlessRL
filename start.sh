# defaultKeepAliveTime  单位min
# defaultPreWarmTime    单位min
# defaultMemoryCapcity  单位GB
# ArricalCnt
# IsFixed 0/1
# SumLimit 
# leftBound
# leftBound2
# rightBound

cleanup() {
    echo "Caught SIGINT signal, terminating the background processes..."
    pkill -P $$
}
trap cleanup SIGINT

keepAliveList=(5 10 30 60 120)

# policyList=('lru' 'lfu' 'mru' 'random' 'maxmem' 'maxmem2' 'maxUsage' 'maxColdStart' 'minColdStart' 'score' 'score3')
policyList=('maxmem')
policyList2=('maxmem' 'cv' 'score1' 'score9' 'score8' 'score7')
# policyList2=('lru' 'lfu' 'random' 'maxmem' 'maxUsage' 'minColdStart' 'score' 'score1' 'score2' 'score3' 'score4')
# policyList2=('lru' 'lfu' 'mru' 'random' 'maxmem' 'maxmem2' 'maxUsage' 'maxColdStart' 'minColdStart' 'score' 'score1' 'score2' 'score3')

memoryList=(500 1000 1500)
arrivalCnt=1

cd pkg/system && go build

for keepAlive in "${keepAliveList[@]}"
do
    for policy in "${policyList[@]}"
    do 
        for memory in "${memoryList[@]}"
        do 
            fixed=1
            prewarm=0
            file="fixed/$policy/fixed-$policy-$keepAlive-$memory-$arrivalCnt"
            echo $file
            ./system $keepAlive $prewarm $memory $arrivalCnt $fixed 0 0 0 0 $policy > ../output/$file.log &
        done
        # wait
    done
    # wait
done
wait
for keepAlive in "${keepAliveList[@]}"
do
    for policy in "${policyList2[@]}"
    do 
        for memory in "${memoryList[@]}"
        do 
            fixed=0
            prewarm=0
            sum=30
            leftBound=0.05
            leftBound2=0.15
            rightBound=0.95
            file="histogram/$policy/histogram-$policy-$keepAlive-$memory-$arrivalCnt"
            echo $file
            ./system $keepAlive $prewarm $memory $arrivalCnt $fixed $sum $leftBound $leftBound2 $rightBound $policy > ../output/$file.log &
        done
        # wait
    done
    # wait
done

wait

# 所有后台进程完成后执行的命令
cd ../.. && ./draw.sh