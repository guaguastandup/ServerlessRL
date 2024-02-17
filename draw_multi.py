import pandas as pd
import matplotlib.pyplot as plt
import sys
import numpy as np

mem_running_usage = {}
mem_occupy_usage = {}
mem_score = {}
time_score = {}
warm_start_rate = {}
cdf_warmstart = {}
app_mem_score = {}
app_time_score = {}

logPath = [
    '0min-1',
    '1s-1',
    '1min-1',
    '5min-1',
    '10min-1',
    '15min-1',
    '20min-1',
    '30min-1', 
    '45min-1',
    '60min-1',
    '90min-1',
    '120min-1'
]

for i in range(len(logPath)):
    with open('./pkg/output/' + logPath[i] + '.log', 'r') as file:
        mem_running_usage[logPath[i]] = []
        mem_occupy_usage[logPath[i]] = []
        mem_score[logPath[i]] = []
        time_score[logPath[i]] = []
        warm_start_rate[logPath[i]] = []
        cdf_warmstart[logPath[i]] = []
        app_mem_score[logPath[i]] = []
        app_time_score[logPath[i]] = []
        for line in file:
            if 'Inf' in line or 'NaN' in line:
                continue
            if line.startswith('MEMRunningUsage'):
                mem_running_usage[logPath[i]].append(float(line.split(': ')[1].split(' ')[0]))
            elif line.startswith('MemOccupyingUsage'):
                mem_occupy_usage[logPath[i]].append(float(line.split(': ')[1].split(' ')[0]))
            elif line.startswith('Mem Score'):
                mem_score[logPath[i]].append(float(line.split(': ')[1].split(' ')[0]))
            elif line.startswith('Time Score'):
                time_score[logPath[i]].append(float(line.split(': ')[1].split(' ')[0]))
            elif line.startswith('warmStart Rate'):
                warm_start_rate[logPath[i]].append(float(line.split(': ')[1].split(' ')[0]))
            elif line.startswith('warmstart rate'):
                cdf_warmstart[logPath[i]].append(float(line.split(':  ')[1].split(' ')[0]))
            elif line.startswith('app mem socre'):
                app_mem_score[logPath[i]].append(float(line.split(': ')[1].split(' ')[0]))
            elif line.startswith('app time socre'):
                app_time_score[logPath[i]].append(float(line.split(': ')[1].split(' ')[0]))

mem_running_usage[logPath[i]] = mem_running_usage[logPath[i]][10: -1]
mem_occupy_usage[logPath[i]] = mem_occupy_usage[logPath[i]][10: -1]
mem_score[logPath[i]] = mem_score[logPath[i]][10: -1]
time_score[logPath[i]] = time_score[logPath[i]][10: -1]
warm_start_rate[logPath[i]] = warm_start_rate[logPath[i]][10: -1]
cdf_warmstart[logPath[i]] = cdf_warmstart[logPath[i]][10: -1]
app_mem_score[logPath[i]] = app_mem_score[logPath[i]][10: -1]
app_time_score[logPath[i]] = app_time_score[logPath[i]][10: -1]

# set plot fontsize
# plt.rcParams.update({'font.size': 20})
plt.rcParams.update({'font.size': 15})

plt.figure(figsize=(25, 25))
plt.subplot(4, 2, 1)
for i in range(len(logPath)):
    plt.plot(mem_running_usage[logPath[i]], label=logPath[i])
plt.legend(loc='upper right')
plt.title('MEMRunningUsage')
plt.xlabel('Minute')
plt.ylabel('Usage (GB)')
plt.grid(True)

plt.subplot(4, 2, 2)
for i in range(len(logPath)):
    plt.plot(mem_occupy_usage[logPath[i]], label=logPath[i])
plt.legend(loc='upper right')
plt.title('MEMOccupyUsage')
plt.xlabel('Minute')
plt.ylabel('Usage (GB)')
plt.grid(True)

plt.subplot(4, 2, 3)
for i in range(len(logPath)):
    plt.plot(mem_score[logPath[i]], label=logPath[i])
plt.legend(loc='upper right')
plt.title('Mem Score')
plt.xlabel('Minute')
plt.ylabel('Score (%)')
plt.grid(True)

plt.subplot(4, 2, 4)
for i in range(len(logPath)):
    plt.plot(time_score[logPath[i]], label=logPath[i])
plt.legend(loc='upper right')
plt.title('Time Score')
plt.xlabel('Minute')
plt.ylabel('Score (%)')
plt.grid(True)

plt.subplot(4, 2, 5)
for i in range(len(logPath)):
    for j in range(len(cdf_warmstart[logPath[i]])):
        cdf_warmstart[logPath[i]][j] = (1.0 - cdf_warmstart[logPath[i]][j]) * 100
    data_sorted = np.sort(cdf_warmstart[logPath[i]])
    cdf = np.arange(1, len(data_sorted) + 1) / len(data_sorted)
    plt.plot(data_sorted, cdf, label=logPath[i])
plt.legend(loc='upper right')
plt.title('CDF of ColdStart Rate')
plt.xlabel('ColdStart Rate (%)')
plt.ylabel('CDF')
plt.grid(True)

plt.subplot(4, 2, 7)
for i in range(len(logPath)):
    data_sorted = np.sort(app_mem_score[logPath[i]])
    cdf = np.arange(1, len(data_sorted) + 1) / len(data_sorted)
    plt.plot(data_sorted, cdf, label=logPath[i])
plt.legend(loc='upper right')
plt.title('CDF of Mem Score')
plt.xlabel('Mem Score (%)')
plt.ylabel('CDF')

plt.subplot(4, 2, 8)
for i in range(len(logPath)):
    data_sorted = np.sort(app_time_score[logPath[i]])
    cdf = np.arange(1, len(data_sorted) + 1) / len(data_sorted)
    plt.plot(data_sorted, cdf, label=logPath[i])
plt.legend(loc='upper right')
plt.title('CDF of Time Score')
plt.xlabel('Time Score (%)')
plt.ylabel('CDF')

plt.tight_layout()
plt.savefig('./pkg/result/multi.png')
plt.close()
