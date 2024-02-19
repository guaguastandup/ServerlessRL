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

# file="5-0-3000-50-fixed"
# ./system 5 0 3000 50 1 0 0 0 0 > ../output/$file.log &

# file="10-0-3000-50-fixed"
# ./system 10 0 3000 50 1 0 0 0 0 > ../output/$file.log &

# file="15-0-3000-50-fixed"
# ./system 15 0 3000 50 1 0 0 0 0 > ../output/$file.log &

# file="30-0-3000-50-fixed"
# ./system 30 0 3000 50 1 0 0 0 0 > ../output/$file.log &

# file="60-0-3000-50-fixed"
# ./system 60 0 3000 50 1 0 0 0 0 > ../output/$file.log &

# file="120-0-3000-50-fixed"
# ./system 120 0 3000 50 1 0 0 0 0 > ../output/$file.log &

# file="5-0-3000-50-histogram"
# ./system 120 0 3000 50 1 50 0.05 0.10 0.95 > ../output/$file.log &

# file="5-0-2000-100-histogram"
# ./system 120 0 2000 50 1 50 0.05 0.10 0.95 > ../output/$file.log

# file="5-0-1500-100-histogram"
# ./system 120 0 1500 50 1 50 0.05 0.10 0.95 > ../output/$file.log

# file="5-0-1000-100-histogram"
# ./system 120 0 1000 50 1 50 0.05 0.10 0.95 > ../output/$file.log

# file="5-0-800-100-histogram"
# ./system 120 0 800 50 1 50 0.05 0.10 0.95 > ../output/$file.log

file="5-0-1500-100-random"
./system 120 0 1500 50 1 50 0.05 0.10 0.95 'random' > ../output/$file.log

file="5-0-1500-100-maxmem"
./system 120 0 1500 50 1 50 0.05 0.10 0.95 'maxmem' > ../output/$file.log

file="5-0-1500-100-maxKeepAlive"
./system 120 0 1500 50 1 50 0.05 0.10 0.95 'maxKeepAlive' > ../output/$file.log

file="5-0-1500-100-minUsage"
./system 120 0 1500 50 1 50 0.05 0.10 0.95 'minUsage' > ../output/$file.log

file="5-0-1500-100-maxColdStartRate"
./system 120 0 1500 50 1 50 0.05 0.10 0.95 'minUsage' > ../output/$file.log

file="5-0-1500-100-fixed-lru" 
./system 120 0 1500 50 1 50 0 0 0 0 'lru' > ../output/$file.log

# 确认当上面所有程序结束后, 再执行以下命令, 该怎么确认？
wait ${pids[@]}

# 所有后台进程完成后执行的命令
cd ../.. && ./draw.sh