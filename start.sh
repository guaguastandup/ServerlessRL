file=$1min-$2
cd pkg/system && go build && ./system $1 > ../output/$file.log 
cd ../.. && python3 draw.py $file