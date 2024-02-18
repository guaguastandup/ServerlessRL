import pandas as pd

for i in range(6, 10):
    mem_file = './app_memory_percentiles.anon.d0' + str(i) + '.csv'
    print(mem_file)
    df = pd.read_csv(mem_file)
    columns_to_keep = ['HashApp', 'AverageAllocatedMb']
    df = df[columns_to_keep]
    df.to_csv('./mem_d0' + str(i) + '.csv', index=False)

for i in range(10, 13):
    invocation_file = './app_memory_percentiles.anon.d' + str(i) + '.csv'
    print(invocation_file)
    df = pd.read_csv(invocation_file)
    columns_to_keep = ['HashApp', 'AverageAllocatedMb']
    df = df[columns_to_keep]
    df.to_csv('./mem_d' + str(i) + '.csv', index=False)