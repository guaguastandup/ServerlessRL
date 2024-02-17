import pandas as pd

prefix = './dataset/azurefunctions-dataset2019/'

for i in range(3, 6):
    mem_file = prefix + 'app_memory_percentiles.anon.d0' + str(i) + '.csv'
    df = pd.read_csv(mem_file) 
    columns_to_keep = ['HashApp', 'AverageAllocatedMb']
    df = df[columns_to_keep]
    df.to_csv('./dataset/azurefunctions/mem_d0' + str(i) + '.csv', index=False)

for i in range(3, 6):
    duration_file = prefix + 'function_durations_percentiles.anon.d0' + str(i) + '.csv'
    df = pd.read_csv(duration_file)
    columns_to_keep = ['HashApp', 'HashFunction', 'Average']
    df = df[columns_to_keep]
    df.to_csv('./dataset/azurefunctions/duration_d0' + str(i) + '.csv', index=False)
    
for i in range(3, 6):
    invacation_file = prefix + 'invocations_per_function_md.anon.d0' + str(i) + '.csv'
    df = pd.read_csv(invacation_file)
    columns_to_drop = ['HashOwner']  # 根据需要替换列名
    df.drop(columns_to_drop, axis=1, inplace=True)
    df.to_csv('./dataset/azurefunctions/invocation_d0' + str(i) + '.csv', index=False)
