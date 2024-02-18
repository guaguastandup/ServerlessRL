file=$1min-$2min-$3GB
cd pkg/system && go build && ./system $1 $2 $3 > ../output/$file.log 
cd ../.. && python3 draw.py $file
./draw.sh