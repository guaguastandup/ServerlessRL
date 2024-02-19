# 如果有第四个参数, 则使用
# $1: 开始时间
# $2: 结束时间
# $3: 内存大小
# $4: 文件名
# 否则使用
if [ ! -n "$4" ]; then
    file=$1-$2m-$3G
else
    file=$1-$2m-$3G-$4
fi
cd pkg/system && go build && ./system $1 $2 $3 > ../output/$file.log 
# cd ../.. && python3 draw.py $file
cd ../.. && ./draw.sh