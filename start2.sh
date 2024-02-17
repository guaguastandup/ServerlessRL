file='KeepAlive-10min-day1-2'
cd pkg/system && go build && ./system 10 > ../output/$file.log 
cd ../.. && python3 draw.py $file