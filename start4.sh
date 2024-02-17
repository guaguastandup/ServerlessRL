file='KeepAlive-20min-day1-2'
cd pkg/system && go build && ./system 20 > ../output/$file.log 
cd ../.. && python3 draw.py $file