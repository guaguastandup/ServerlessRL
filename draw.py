import pandas as pd
import matplotlib.pyplot as plt
import sys
import numpy as np

# 初始化数据存储列表
mem_running_usage = []
mem_score = []
time_score = []
warm_start_rate = []
cdf_warmstart = []

logPath = sys.argv[1]

# 打开并逐行读取a.log文件
with open('./pkg/output/' + logPath + '.log', 'r') as file: 
    for line in file:
        if 'Inf' in line or 'NaN' in line:
            continue
        # 对每一行进行解析，提取感兴趣的数据
        if line.startswith('MEMRunningUsage'):
            mem_running_usage.append(float(line.split(': ')[1].split(' ')[0]))
        elif line.startswith('Mem Score'):
            mem_score.append(float(line.split(': ')[1].split(' ')[0]))
        elif line.startswith('Time Score'):
            time_score.append(float(line.split(': ')[1].split(' ')[0]))
        elif line.startswith('warmStart Rate'):
            warm_start_rate.append(float(line.split(': ')[1].split(' ')[0]))
        elif line.startswith('warmstart rate'):
            cdf_warmstart.append(float(line.split(':  ')[1].split(' ')[0]))

# 创建DataFrame
mem_running_usage = mem_running_usage[20:]
mem_score = mem_score[20:]
time_score = time_score[20:]
warm_start_rate = warm_start_rate[20:]

# 绘图
# plt.figure(figsize=(15, 8))
plt.figure(figsize=(12, 6))
# 绘制MEMRunningUsage
plt.subplot(2, 2, 1)
plt.plot(mem_running_usage, linestyle='-', linewidth=2)
# plt.plot(df['MEMRunningUsage (GB)'], marker='o', linestyle='-')
plt.title('MEMRunningUsage')
plt.xlabel('Sample')
plt.ylabel('Usage (GB)')

# 绘制Mem Score
plt.subplot(2, 2, 2)
plt.plot(mem_score, linestyle='-', linewidth=2, color='green')
# plt.plot(df['Mem Score (%)'], marker='o', linestyle='-', color='green')
plt.title('Mem Score')
plt.xlabel('Sample')
plt.ylabel('Score (%)')

# 绘制Time Score
plt.subplot(2, 2, 3)
plt.plot(time_score, linestyle='-', linewidth=2, color='red')
# plt.plot(df['Time Score (%)'], marker='o', linestyle='-', color='red')
plt.title('Time Score')
plt.xlabel('Sample')
plt.ylabel('Score (%)')

# 绘制WarmStart Rate
plt.subplot(2, 2, 4)
plt.plot(warm_start_rate, linestyle='-', linewidth=2, color='purple')
# plt.plot(df['WarmStart Rate (%)'], marker='o', linestyle='-', color='purple')
plt.title('WarmStart Rate')
plt.xlabel('Sample')
plt.ylabel('Rate (%)')

plt.tight_layout()
plt.savefig('./pkg/output/' + logPath + '.png')
plt.close()

# 绘制CDF
plt.figure(figsize=(12, 6))
# cdf_warmstart.sort()
# yvals = range(1, len(cdf_warmstart) + 1) / len(cdf_warmstart)
# plt.plot(cdf_warmstart, yvals, linestyle='-', linewidth=2, color='purple')
# plt.plot(cdf_warmstart, yvals, marker='o', linestyle='-', color='purple')

# 计算CDF
# data_sorted = np.array([1, 1, 2, 5, 5])

data_sorted = np.sort(cdf_warmstart)
cdf = np.arange(1, len(data_sorted) + 1) / len(data_sorted)
# 绘制CDF图
plt.plot(data_sorted, cdf, marker='.', linestyle='-')
# plt.plot(cdf, data_sorted, marker='.', linestyle='-')
plt.title('CDF of WarmStart Rate')
plt.xlabel('Sample')
plt.ylabel('Rate (%)')
plt.savefig('./pkg/output/' + logPath + '_warm_cdf.png')
plt.close()

for i in range(len(cdf_warmstart)):
    cdf_warmstart[i] = 1.0 - cdf_warmstart[i]
data_sorted = np.sort(cdf_warmstart)
cdf = np.arange(1, len(data_sorted) + 1) / len(data_sorted)
# 绘制CDF图
plt.plot(data_sorted, cdf, marker='.', linestyle='-')
# plt.plot(cdf, data_sorted, marker='.', linestyle='-')
plt.title('CDF of ColdStart Rate')
plt.xlabel('Sample')
plt.ylabel('Rate (%)')
plt.savefig('./pkg/output/' + logPath + '_cold_cdf.png')
plt.close()
