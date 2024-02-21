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

# go build && ./system 5 0 1500 50 0 50 0.05 0.1 0.95 lru > ../output/try.log

keepAliveList=(5)
# policyList=('random' 'lru' 'maxmem' 'maxKeepAlive' 'minUsage' 'maxColdStartRate')
policyList=('score' 'score1' 'score2' 'score3' 'maxKeepAlive' 'minKeepAlive' 'maxUsage' 'minColdStartRate')
# memoryList=(500 900 1300 1700)
memoryList=(500 1300)

cd pkg/system && go build

for keepAlive in "${keepAliveList[@]}"
do
    for policy in "${policyList[@]}"
    do 
        for memory in "${memoryList[@]}"
        do 
            fixed=1
            arrivalCnt=50
            prewarm=0
            file="fixed/$policy/fixed-$policy-$keepAlive-$prewarm-$memory-$arrivalCnt"
            echo $file
            ./system $keepAlive $prewarm $memory $arrivalCnt $fixed 0 0 0 0 $policy > ../output/$file.log &
        done
        # wait
    done
    wait
done

wait

for keepAlive in "${keepAliveList[@]}"
do
    for policy in "${policyList[@]}"
    do 
        for memory in "${memoryList[@]}"
        do 
            fixed=0
            arrivalCnt=50
            prewarm=0
            sum=50
            leftBound=0.05
            leftBound2=0.10
            rightBound=0.95
            file="histogram/$policy/histogram-$policy-$keepAlive-$prewarm-$memory-$arrivalCnt"
            echo $file
            ./system $keepAlive $prewarm $memory $arrivalCnt $fixed $sum $leftBound $leftBound2 $rightBound $policy > ../output/$file.log &
        done
        # wait
    done
    wait
done

wait

# 所有后台进程完成后执行的命令
# cd ../.. && ./draw.sh