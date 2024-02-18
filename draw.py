import pandas as pd
import matplotlib.pyplot as plt
import sys
import numpy as np

mem_running_usage = []
mem_occupy_usage = []
mem_score = []
time_score = []
warm_start_rate = []
cdf_warmstart = []
app_mem_score = []
app_time_score = []

logPath = sys.argv[1]

# 打开并逐行读取a.log文件
with open('./pkg/output/' + logPath + '.log', 'r') as file: 
    for line in file:
        if 'Inf' in line or 'NaN' in line:
            continue
        # 对每一行进行解析，提取感兴趣的数据
        if line.startswith('MEMRunningUsage'):
            mem_running_usage.append(float(line.split(': ')[1].split(' ')[0]))
        elif line.startswith('MemOccupyingUsage'):
            mem_occupy_usage.append(float(line.split(': ')[1].split(' ')[0]))
        elif line.startswith('Mem Score'):
            mem_score.append(float(line.split(': ')[1].split(' ')[0]))
        elif line.startswith('Time Score'):
            time_score.append(float(line.split(': ')[1].split(' ')[0]))
        elif line.startswith('warmStart Rate'):
            warm_start_rate.append(float(line.split(': ')[1].split(' ')[0]))
        elif line.startswith('warmstart rate'):
            cdf_warmstart.append(float(line.split(':  ')[1].split(' ')[0]))
        elif line.startswith('app mem socre'):
            app_mem_score.append(float(line.split(': ')[1].split(' ')[0]))
        elif line.startswith('app time socre'):
            app_time_score.append(float(line.split(': ')[1].split(' ')[0]))


# 创建DataFrame
# 删除第一个和最后一个元素
mem_running_usage = mem_running_usage[10: -5]
mem_occupy_usage = mem_occupy_usage[10: -5]
mem_score = mem_score[10: -5]
time_score = time_score[10: -5]
warm_start_rate = warm_start_rate[10: -5]
cdf_warmstart = cdf_warmstart[10: -5]
app_mem_score = app_mem_score[10: -5]
app_time_score = app_time_score[10: -5]

plt.figure(figsize=(15, 10))
plt.rcParams.update({'font.size': 15})

# # 绘制MEMRunningUsage
plt.subplot(2, 2, 1)
plt.plot(mem_running_usage, label='Memory Running', color='red')
plt.plot(mem_occupy_usage, label='Memory Occupying', color='blue')
plt.legend(loc='upper right')
plt.title('MEMRunningUsage and MEMOccupyUsage')
plt.xlabel('Minute')
plt.ylabel('Usage (GB)')
plt.grid(True)

# 绘制Mem Score & Time Score
plt.subplot(2, 2, 2)
plt.plot(mem_score, label='Mem Score', color='green')
plt.plot(time_score, label='Time Score', color='orange')
plt.legend(loc='upper right')
plt.title('Mem Score & Time Score')
plt.xlabel('Minute')
plt.ylabel('Score (%)')
plt.grid(True)

# 绘制冷启动CDF图
plt.subplot(2, 2, 3)
for i in range(len(cdf_warmstart)):
    cdf_warmstart[i] = (1.0 - cdf_warmstart[i]) * 100
data_sorted = np.sort(cdf_warmstart)
cdf = np.arange(1, len(data_sorted) + 1) / len(data_sorted)
plt.plot(data_sorted, cdf, color='purple')
plt.title('CDF of ColdStart Rate')
plt.xlabel('ColdStart Rate (%)')
plt.ylabel('CDF')

# 绘制冷启动CDF图
plt.subplot(2, 2, 4)
data_sorted = np.sort(app_mem_score)
cdf = np.arange(1, len(data_sorted) + 1) / len(data_sorted)
plt.plot(data_sorted, cdf, color='pink', label='Mem Score')
data_sorted = np.sort(app_time_score)
cdf = np.arange(1, len(data_sorted) + 1) / len(data_sorted)
plt.plot(data_sorted, cdf, color='grey', label='Time Score')
plt.legend(loc='upper right')
plt.title('CDF of Score')
plt.xlabel('Score (%)')
plt.ylabel('CDF')

plt.tight_layout()
plt.savefig('./pkg/result/' + logPath + '.png')
plt.close()
