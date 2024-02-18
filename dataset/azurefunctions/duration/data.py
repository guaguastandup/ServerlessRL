import pandas as pd

for i in range(6, 10):
    duration_file = './function_durations_percentiles.anon.d0' + str(i) + '.csv'
    print(duration_file)
    df = pd.read_csv(duration_file)
    columns_to_keep = ['HashApp', 'HashFunction', 'Average']
    df = df[columns_to_keep]
    df.to_csv('./duration_d0' + str(i) + '.csv', index=False)

for i in range(10, 13):
    duration_file = './function_durations_percentiles.anon.d' + str(i) + '.csv'
    print(duration_file)
    df = pd.read_csv(duration_file)
    columns_to_keep = ['HashApp', 'HashFunction', 'Average']
    df = df[columns_to_keep]
    df.to_csv('./duration_d' + str(i) + '.csv', index=False)