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
    # 杀死所有属于当前脚本进程组的后台进程
    pkill -P $$
}

trap cleanup SIGINT


cd pkg/system && go build

for keepAlive in 5 10 15 30 60 120
do
    for policy in 'random' 'maxmem' 'maxKeepAlive' 'minUsage' 'maxColdStartRate' 'lru'
    do 
        fixed=1
        memory=3000
        arrivalCnt=50
        prewarm=0
        policy='lru'
        file="fixed-$policy-$i-$prewarm-$memory-$arrivalCnt"
        ./system $keepAlive $prewarm $memory $arrivalCnt $fixed 0 0 0 0 $policy > ../output/$file.log &
    done
done

for keepAlive in 5 10 15 30 60 120
do
    for policy in 'random' 'maxmem' 'maxKeepAlive' 'minUsage' 'maxColdStartRate' 'lru'
    do 
        fixed=0
        memory=3000
        arrivalCnt=50
        prewarm=0
        policy='lru'
        file="fixed-$policy-$i-$prewarm-$memory-$arrivalCnt"
        sum=50
        leftBound=0.05
        leftBound2=0.10
        rightBound=0.95
        ./system $keepAlive $prewarm $memory $arrivalCnt $fixed $sum $leftBound $leftBound2 $rightBound $policy > ../output/$file.log &
    done
done

wait ${pids[@]}

# 所有后台进程完成后执行的命令
cd ../.. && ./draw.sh