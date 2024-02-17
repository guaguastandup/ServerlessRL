# cd pkg/system_funconly && go build && ./system_funconly > ../output/$1.log 
cd pkg/system && go build && ./system > ../output/$1.log 
cd ../.. && python3 draw.py $1