file=$1min-$2min
cd pkg/system && go build && ./system $1 $2 > ../output/$file.log 
cd ../.. && python3 draw.py $file
./draw.sh