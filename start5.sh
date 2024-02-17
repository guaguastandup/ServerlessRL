file='KeepAlive-100min-day1-2'
cd pkg/system && go build && ./system 100 > ../output/$file.log 
cd ../.. && python3 draw.py $file