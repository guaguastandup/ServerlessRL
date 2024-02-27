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

keepAliveList=(10)
# policyList=('lru' 'lfu' 'maxmem' 'random' 'score1')
# policyList2=('lru' 'lfu' 'maxmem' 'random' 'score1')
policyList=('maxmem' 'score1')
policyList2=('maxmem' 'score1')
# memoryList=(200 400 600 800 1000)
memoryList=(400 800 1200 1600)
arrivalCnt=1

cd pkg/system && go build
for keepAlive in "${keepAliveList[@]}"
do
    for memory in "${memoryList[@]}"
    do 
        for policy in "${policyList[@]}"
        do 
            fixed=1
            prewarm=0
            # file name
            dir="../output/fixed/$policy" # 设置目录路径
            file="fixed-$policy-$keepAlive-$memory-$arrivalCnt.log" # 设置文件名
            fullpath="$dir/$file" # 完整的文件路径
            mkdir -p "$dir"
            echo "$fullpath"
            ./system $keepAlive $prewarm $memory $arrivalCnt $fixed 0 0 0 0 $policy > ../output/$fullpath &
        done
        # wait
    done
    # wait
done
# wait
for keepAlive in "${keepAliveList[@]}"
do
    for memory in "${memoryList[@]}"
    do
        for policy in "${policyList2[@]}"
        do  
            fixed=0
            prewarm=0
            sum=20
            leftBound=0.05
            leftBound2=0.15
            rightBound=0.95
            # file name
            dir="../output/histogram/$policy" # 设置目录路径
            file="histogram-$policy-$keepAlive-$memory-$arrivalCnt.log" # 设置文件名
            fullpath="$dir/$file" # 完整的文件路径
            mkdir -p "$dir"
            echo "$fullpath"
            ./system $keepAlive $prewarm $memory $arrivalCnt $fixed $sum $leftBound $leftBound2 $rightBound $policy > ../output/$fullpath &
        done
        # wait
    done
    # wait
done
wait
cd ../.. && ./draw.sh