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



# ----------------------  fixed ------------------------------
# interate
# file="fixed-lru-5-0-3000-50"
# ./system 5 0 3000 50 1 0 0 0 0 > ../output/$file.log &

# file="fixed-lru-10-0-3000-50"
# ./system 10 0 3000 50 1 0 0 0 0 > ../output/$file.log &

# ------------------  fixed & policy ------------------------------
# file="fixed-lru-5-0-3000-50"
# ./system 5 0 3000 50 1 0 0 0 0 > ../output/$file.log &

# file="fixed-lru-10-0-3000-50"
# ./system 10 0 3000 50 1 0 0 0 0 > ../output/$file.log &

# ---------------------- histogram -----------------------------
# file="histogram-none-5-0-3000"
# ./system 5 0 3000 50 1 50 0.05 0.10 0.95 lru > ../output/$file.log &

# ------------------ histogram & policy-----------------------------
file="5-0-600-50-maxmem"
./system 5 0 600 50 1 50 0.05 0.10 0.95 maxmem > ../output/$file.log &
file="5-0-800-50-maxmem"
./system 5 0 800 50 1 50 0.05 0.10 0.95 maxmem > ../output/$file.log &
file="5-0-1000-50-maxmem"
./system 5 0 1000 50 1 50 0.05 0.10 0.95 maxmem > ../output/$file.log &

# file="5-0-1500-50-random"
# ./system 5 0 1500 50 1 50 0.05 0.10 0.95 random > ../output/$file.log &

# file="5-0-1500-50-maxKeepAlive"
# ./system 5 0 1500 50 1 50 0.05 0.10 0.95 maxKeepAlive > ../output/$file.log &

# file="5-0-1500-50-minUsage"
# ./system 5 0 1500 50 1 50 0.05 0.10 0.95 minUsage > ../output/$file.log &

# file="5-0-1500-50-maxColdStartRate"
# ./system 5 0 1500 50 1 50 0.05 0.10 0.95 maxColdStartRate > ../output/$file.log &

# file="5-0-1500-50-fixed-lru" 
# ./system 5 0 1500 50 1 50 0 0 0 lru > ../output/$file.log &

# file="5-0-1500-50-fixed-random" 
# ./system 5 0 1500 50 1 50 0 0 0 random > ../output/$file.log &

# file="30-0-1500-50-fixed-random" 
# ./system 30 0 1500 50 1 50 0 0 0 lru > ../output/$file.log &

# file="30-0-1500-50-fixed-maxmem" 
# ./system 30 0 1500 50 1 50 0 0 0 maxmem > ../output/$file.log &

# 确认当上面所有程序结束后, 再执行以下命令, 该怎么确认？
wait ${pids[@]}

# 所有后台进程完成后执行的命令
cd ../.. && ./draw.sh