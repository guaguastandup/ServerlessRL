import pandas as pd
import matplotlib.pyplot as plt
import sys
import numpy as np
mem_running_usage, mem_occupy_usage = {}, {}
mem_score, time_score = {}, {}
warm_start_rate, cdf_warmstart = {}, {}
app_mem_score, app_time_score = {}, {}
avg_coldstart, avg_mem_score, avg_time_score = {}, {}, {}
warm_start, timeCost, memCost = {}, {}, {}
keepAliveList=[15]
policyList=[
    # 'random',
    # 'lru',
    # 'maxmem', 
    # 'score1',
    # 'score2',
    # 'score3',
    # 'score4',
    # 'score5',
    # "ideal",
]

policyList2=[
    'random',
    'lru',
    'maxmem', 
    # 'score1',
    # 'score2',
    # 'score3',
    'score4',
    # 'score5',
    # "ideal",
]
# memoryList=[200, 400, 600]
memoryList=[600]
arrivalCnt=1
id = 0
maxlen = 520

def draw_cdf():
    plt.rcParams.update({'font.size': 20})
    plt.figure(figsize=(10, 8))
    for i in range(len(logPath)):
        for j in range(len(cdf_warmstart[logPath[i]])):
            cdf_warmstart[logPath[i]][j] = 100.0 * (1.0 - cdf_warmstart[logPath[i]][j])
    for i in range(len(logPath)):
        data_sorted = np.sort(cdf_warmstart[logPath[i]])
        cdf = np.arange(1, len(data_sorted) + 1) / len(data_sorted)
        plt.plot(data_sorted, cdf, label=label[logPath[i]])
    plt.legend(loc='upper right')
    # plt.title('CDF of ColdStart Rate')
    plt.xlabel('ColdStart Rate (%)')
    plt.ylabel('CDF')
    plt.grid(True)
    plt.savefig(f"./pkg/result/cdf_multi-{str(id)}.pdf")
    plt.close()

def statistics():
    for file in logPath:
        dot20, dot40, dot60, dot80, dot99, dot100 = 0, 0, 0, 0, 0, 0
        for i in range(len(cdf_warmstart[file])):
            if 0 <= cdf_warmstart[file][i] < 20:
                dot20 += 1
            elif 20 <= cdf_warmstart[file][i] < 40:
                dot40 += 1
            elif 40 <= cdf_warmstart[file][i] < 60:
                dot60 += 1
            elif 60 <= cdf_warmstart[file][i] < 80:
                dot80 += 1
            elif 80 <= cdf_warmstart[file][i] < 100:
                dot99 += 1
            elif 100 == cdf_warmstart[file][i]:
                dot100 += 1
        avg = sum(cdf_warmstart[file]) / len(cdf_warmstart[file])
        print(label[file], avg)
        print(label[file], dot20, dot40, dot60, dot80, dot99, dot100, '\n')

def read():
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
            avg_coldstart[logPath[i]] = []
            avg_mem_score[logPath[i]] = []
            avg_time_score[logPath[i]] = []
            warm_start[logPath[i]] = []
            memCost[logPath[i]] = []
            timeCost[logPath[i]] = []
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
                elif line.startswith('app average coldstart rate'):
                    avg_coldstart[logPath[i]].append(float(line.split(': ')[1].split(' ')[0]))
                elif line.startswith('average mem socre'):
                    avg_mem_score[logPath[i]].append(float(line.split(': ')[1].split(' ')[0]))
                elif line.startswith('average time socre'):
                    avg_time_score[logPath[i]].append(float(line.split(': ')[1].split(' ')[0]))
                elif line.startswith('warmstart'):
                    warm_start[logPath[i]].append(float(line.split(': ')[1].split(' ')[0]))
                elif line.startswith('total coldstartTimeCost'):
                    timeCost[logPath[i]].append(float(line.split(': ')[1].split(' ')[0]))
                elif line.startswith('total coldstartMemCost'):
                    memCost[logPath[i]].append(float(line.split(': ')[1].split(' ')[0]))
for i in range(len(memoryList)):
    logPath = []
    label = {}
    id = memoryList[i]
    for k in keepAliveList[0:1]:
        for p in policyList:
            for m in memoryList[i:i+1]:
                if "ideal" in p:
                    logPath.append('fixed/' + p + '/fixed-' + p + '-' + str(k) + '-' + str(100000) + '-' + str(arrivalCnt))
                else:
                    logPath.append('fixed/' + p + '/fixed-' + p + '-' + str(k) + '-' +str(m) + '-' + str(arrivalCnt))
    for k in keepAliveList[0:1]:
        for p in policyList2:
            for m in memoryList[i:i+1]:
                if "ideal" in p:
                    logPath.append('histogram/' + p + '/histogram-' + p + '-' + str(k) + '-' + str(100000) + '-' + str(arrivalCnt))
                else:
                    logPath.append('histogram/' + p + '/histogram-' + p + '-' + str(k) + '-' + str(m) + '-' + str(arrivalCnt))
    
    for path in logPath:
        # type = path.split('/')[0]
        # policy = path.split('/')[1]
        # keepAlive = path.split('-')[2]
        # memory = path.split('-')[3]
        # label[path] = type + '-' + policy + '-' + keepAlive + '-' + memory
        if 'maxmem' in path:
            label[path] = 'MaxMem'
        elif 'lru' in path:
            label[path] = 'LRU'
        elif 'random' in path:
            label[path] = 'Random'
        elif 'score4' in path:
            label[path] = 'FISH'
        # elif 'score5' in path:
            # label[path] = 'Memory&Percentage'
    read()
    # draw_score()
    # draw_cost()
    # draw_warmstart()
    draw_cdf()
    statistics()
    print("-----------------------------------------------------------")