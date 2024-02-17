file='KeepAlive-15min-day1-2'
cd pkg/system && go build && ./system 15 > ../output/$file.log 
cd ../.. && python3 draw.py $file