import pandas as pd

for i in range(1, 3):
    region_file = './region_0' + str(i) + '.csv'
    print(region_file)
    df = pd.read_csv(region_file)
    
    columns_to_keep = ['__time__', 'functionName', 'latency', 'memoryMB']
    df = df[columns_to_keep]
    df.to_csv('./data_' + str(i) + '.csv', index=False)