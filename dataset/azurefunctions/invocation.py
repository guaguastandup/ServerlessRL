import pandas as pd

# for i in range(6, 10):
#     invocation_file = './invocations_per_function_md.anon.d0' + str(i) + '.csv'
#     print(invocation_file)
#     df = pd.read_csv(invocation_file)
#     columns_to_drop = ['HashOwner']
#     df = df.drop(columns_to_drop, axis=1)
#     df.to_csv('./invocation_d0' + str(i) + '.csv', index=False)

# for i in range(10, 13):
#     invocation_file = './invocations_per_function_md.anon.d' + str(i) + '.csv'
#     print(invocation_file)
#     df = pd.read_csv(invocation_file)
#     columns_to_drop = ['HashOwner']
#     df = df.drop(columns_to_drop, axis=1)
#     df.to_csv('./invocation_d' + str(i) + '.csv', index=False)
       
# for i in range(3, 4):
#     invocation_file = './invocation_d0' + str(i) + '.csv'
#     df = pd.read_csv(invocation_file) 
#     for j in range(154, 1141):
#         columns_to_keep = ['HashApp', 'HashFunction', 'Trigger', str(j)]
#         df1 = df[columns_to_keep]
#         df1 = df1[df1[str(j)] != 0]
#         df1.to_csv('./invocation/d0' + str(i) + '/invocation_d0' + str(i) + '_m' + str(j) + '.csv', index=False)
        
# for i in range(4, 10):
#     invocation_file = './invocation_d0' + str(i) + '.csv'
#     df = pd.read_csv(invocation_file) 
#     for j in range(1, 1141):
#         columns_to_keep = ['HashApp', 'HashFunction', 'Trigger', str(j)]
#         df1 = df[columns_to_keep]
#         df1 = df1[df1[str(j)] != 0]
#         df1.to_csv('./invocation/d0' + str(i) + '/invocation_d0' + str(i) + '_m' + str(j) + '.csv', index=False)
        
for i in range(10, 13):
    invocation_file = './invocation_d' + str(i) + '.csv'
    df = pd.read_csv(invocation_file) 
    for j in range(1, 1141):
        columns_to_keep = ['HashApp', 'HashFunction', 'Trigger', str(j)]
        df1 = df[columns_to_keep]
        df1 = df1[df1[str(j)] != 0]
        df1.to_csv('./invocation/d' + str(i) + '/invocation_d' + str(i) + '_m' + str(j) + '.csv', index=False)