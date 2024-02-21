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

logPath = []

# keepAliveList=[5, 10, 15, 30, 60, 120]
keepAliveList=[10]
policyList=[
    'random', 
    'lru', 
    'maxmem', 
    'minUsage', 'maxUsage',
    # 'minColdStartRate', 'maxColdStartRate',
    # 'score', 'score1', 'score2', 'score3', 
    # 'maxKeepAlive', 'minKeepAlive',
]
# memoryList=[500, 900, 1300, 1700]
memoryList=[900]

id = 0
if len(sys.argv) > 1:
    id = int(sys.argv[1])
    
for k in keepAliveList[0:1]:
    for p in policyList:
        for m in memoryList[0:1]:
            logPath.append('fixed-' + p + '-' + str(k) + '-0-' + str(m) + '-50')
            logPath.append('histogram-' + p + '-' + str(k) + '-0-' + str(m) + '-50')

def draw():
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

    # mem_running_usage[logPath[i]] = mem_running_usage[logPath[i]][10: -5]
    # mem_occupy_usage[logPath[i]] = mem_occupy_usage[logPath[i]][10: -5]
    # mem_score[logPath[i]] = mem_score[logPath[i]][10: -5]
    # time_score[logPath[i]] = time_score[logPath[i]][10: -5]
    # warm_start_rate[logPath[i]] = warm_start_rate[logPath[i]][10: -5]
    # cdf_warmstart[logPath[i]] = cdf_warmstart[logPath[i]][10: -5]
    # app_mem_score[logPath[i]] = app_mem_score[logPath[i]][10: -5]
    # app_time_score[logPath[i]] = app_time_score[logPath[i]][10: -5]

    # set plot fontsize
    # plt.rcParams.update({'font.size': 20})
    plt.rcParams.update({'font.size': 15})

    plt.figure(figsize=(20, 15))
    plt.subplot(3, 2, 1)
    plt.plot(mem_running_usage[logPath[0]], label='Mem Running')
    for i in range(len(logPath)):
        plt.plot(mem_occupy_usage[logPath[i]], label=logPath[i])
    plt.legend(loc='upper right')
    plt.title('MEMOccupyUsage')
    plt.xlabel('Minute')
    plt.ylabel('Usage (GB)')
    plt.grid(True)

    plt.subplot(3, 2, 2)
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

    plt.subplot(3, 2, 3)
    for i in range(len(logPath)):
        plt.plot(mem_score[logPath[i]], label=logPath[i])
    plt.legend(loc='upper right')
    plt.title('Mem Score')
    plt.xlabel('Minute')
    plt.ylabel('Score (%)')
    plt.grid(True)

    plt.subplot(3, 2, 4)
    for i in range(len(logPath)):
        plt.plot(time_score[logPath[i]], label=logPath[i])
    plt.legend(loc='upper right')
    plt.title('Time Score')
    plt.xlabel('Minute')
    plt.ylabel('Score (%)')
    plt.grid(True)


    plt.subplot(3, 2, 5)
    for i in range(len(logPath)):
        data_sorted = np.sort(app_mem_score[logPath[i]])
        cdf = np.arange(1, len(data_sorted) + 1) / len(data_sorted)
        plt.plot(data_sorted, cdf, label=logPath[i])
    plt.legend(loc='upper right')
    plt.title('CDF of Mem Score')
    plt.xlabel('Mem Score (%)')
    plt.ylabel('CDF')

    plt.subplot(3, 2, 6)
    for i in range(len(logPath)):
        data_sorted = np.sort(app_time_score[logPath[i]])
        cdf = np.arange(1, len(data_sorted) + 1) / len(data_sorted)
        plt.plot(data_sorted, cdf, label=logPath[i])
    plt.legend(loc='upper right')
    plt.title('CDF of Time Score')
    plt.xlabel('Time Score (%)')
    plt.ylabel('CDF')

    plt.tight_layout()
    plt.savefig(f"./pkg/result/multi-{str(id)}.png")
    plt.close()

    plt.figure(figsize=(15, 25))
    for i in range(len(logPath)):
        data_sorted = np.sort(cdf_warmstart[logPath[i]])
        cdf = np.arange(1, len(data_sorted) + 1) / len(data_sorted)
        plt.plot(data_sorted, cdf, label=logPath[i])
    plt.legend(loc='upper right')
    plt.title('CDF of ColdStart Rate')
    plt.xlabel('ColdStart Rate (%)')
    plt.ylabel('CDF')
    plt.grid(True)
    plt.savefig(f"./pkg/result/cdf_multi-{str(id)}.png")
    plt.close()

    for file in logPath:
        dot10 = 0
        dot25 = 0
        dot50 = 0
        dot75 = 0
        dot90 = 0    
        dot100 = 0
        for i in range(len(cdf_warmstart[file])):
            if 0 <= cdf_warmstart[file][i] < 10:
                dot10 += 1
            elif 10 <= cdf_warmstart[file][i] < 25:
                dot25 += 1
            elif 25 <= cdf_warmstart[file][i] < 50:
                dot50 += 1
            elif 50 <= cdf_warmstart[file][i] < 75:
                dot75 += 1
            elif 75 <= cdf_warmstart[file][i] < 90:
                dot90 += 1
            elif 90 <= cdf_warmstart[file][i] <= 100:
                dot100 += 1
        print(file, sum(cdf_warmstart[file]) / len(cdf_warmstart[file]))
        print(file, dot10, dot25, dot50, dot75, dot90, dot100, '\n')
        
    
draw()