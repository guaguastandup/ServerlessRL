file='KeepAlive-5min-day1-3'
cd pkg/system && go build && ./system 5 > ../output/$file.log 
cd ../.. && python3 draw.py $file