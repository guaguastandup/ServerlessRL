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
keepAliveList=[10]
policyList=[
    # 'random',
    # 'lru',
    # 'maxmem', 
    # 'score1',
    # 'score2',
    # 'score3',
    'score4',
    # 'score5',
]
policyList2=[
    # 'random',
    # 'lru',
    # 'maxmem', 
    # 'score1',
    # 'score2',
    # 'score3',
    'score4',
    # 'score5',
]

# memoryList=[200, 400, 600, 800]
memoryList=[500]
arrivalCnt=1
id = 0

# maxlen = 520
# maxlen = 200
maxlen = 0
def draw():
    plt.rcParams.update({'font.size': 15})
    plt.figure(figsize=(40, 20))
    plt.subplot(3, 2, 1)
    plt.plot(mem_running_usage[logPath[0]], label='Mem Running')
    for i in range(len(logPath)):
        plt.plot(mem_occupy_usage[logPath[i]], label=label[logPath[i]])
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
        plt.plot(data_sorted, cdf, label=label[logPath[i]])
    plt.legend(loc='upper right')
    plt.title('CDF of ColdStart Rate')
    plt.xlabel('ColdStart Rate (%)')
    plt.ylabel('CDF')
    plt.grid(True)

    plt.subplot(3, 2, 3)
    for i in range(len(logPath)):
        plt.plot(mem_score[logPath[i]], label=label[logPath[i]])
    plt.legend(loc='upper right')
    plt.title('Mem Score')
    plt.xlabel('Minute')
    plt.ylabel('Score (%)')
    plt.grid(True)

    plt.subplot(3, 2, 4)
    for i in range(len(logPath)):
        plt.plot(time_score[logPath[i]], label=label[logPath[i]])
    plt.legend(loc='upper right')
    plt.title('Time Score')
    plt.xlabel('Minute')
    plt.ylabel('Score (%)')
    plt.grid(True)

    plt.subplot(3, 2, 5)
    for i in range(len(logPath)):
        data_sorted = np.sort(app_mem_score[logPath[i]])
        cdf = np.arange(1, len(data_sorted) + 1) / len(data_sorted)
        plt.plot(data_sorted, cdf, label=label[logPath[i]])
    plt.legend(loc='upper right')
    plt.title('CDF of Mem Score')
    plt.xlabel('Mem Score (%)')
    plt.ylabel('CDF')
    
    plt.subplot(3, 2, 6)
    for i in range(len(logPath)):
        data_sorted = np.sort(app_time_score[logPath[i]])
        cdf = np.arange(1, len(data_sorted) + 1) / len(data_sorted)
        plt.plot(data_sorted, cdf, label=label[logPath[i]])
    plt.legend(loc='upper right')
    plt.title('CDF of Time Score')
    plt.xlabel('Time Score (%)')
    plt.ylabel('CDF')
    
    plt.tight_layout()
    plt.savefig(f"./pkg/result/multi-{str(id)}.png")
    plt.close()

def draw_cdf():
    plt.figure(figsize=(15, 25))
    for i in range(len(logPath)):
        data_sorted = np.sort(cdf_warmstart[logPath[i]])
        cdf = np.arange(1, len(data_sorted) + 1) / len(data_sorted)
        plt.plot(data_sorted, cdf, label=label[logPath[i]])
    plt.legend(loc='upper right')
    plt.title('CDF of ColdStart Rate')
    plt.xlabel('ColdStart Rate (%)')
    plt.ylabel('CDF')
    plt.grid(True)
    plt.savefig(f"./pkg/result/cdf_multi-{str(id)}.png")
    plt.close()

def statistics():
    for file in logPath:
        dot10, dot25, dot50, dot75, dot90, dot100 = 0, 0, 0, 0, 0, 0
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
        avg = sum(cdf_warmstart[file]) / len(cdf_warmstart[file])
        print(label[file], avg)
        print(label[file], dot10, dot25, dot50, dot75, dot90, dot100, '\n')

def draw_score():
    plt.figure(figsize=(15, 25))
    plt.rcParams.update({'font.size': 30})
    minlen = 1000000000
    h_m = 1
    for i in range(len(logPath)):
        plt.plot(avg_coldstart[logPath[i]], label=label[logPath[i]], linewidth=2)
        minlen = min(minlen, len(avg_coldstart[logPath[i]]))
    print("minlen: ", minlen)
    best = ''
    best_val = 100000
    
    if maxlen != 0 :
        minlen = maxlen
    for i in range(len(logPath)):
        if 'his' in logPath[i] and 'maxmem' in logPath[i]:
            h_m = avg_coldstart[logPath[i]][minlen-1]
        if best_val > avg_coldstart[logPath[i]][minlen-1] and 'ideal' not in label[logPath[i]]:
            best = label[logPath[i]]
            best_val = avg_coldstart[logPath[i]][minlen-1]
        print(logPath[i], avg_coldstart[logPath[i]][minlen-1])
    print("\nbest coldstart rate: ", best, best_val, "decrease: ", 100.0 * (h_m - best_val)/h_m, "%")
    plt.legend(loc='upper left')
    plt.title('Average ColdStart Rate')
    plt.xlabel('Minute')
    plt.ylabel('ColdStart Rate')
    plt.grid(True)
    plt.tight_layout()
    plt.savefig(f"./pkg/result/multi-score-{str(id)}.png")
    plt.close()
    
def draw_warmstart():
    plt.figure(figsize=(15, 25))
    plt.rcParams.update({'font.size': 30})
    minlen = 1000000000
    h_m = 1
    for i in range(len(logPath)):
        plt.plot(warm_start[logPath[i]], label=label[logPath[i]], linewidth=2)
        minlen = min(minlen, len(warm_start[logPath[i]]))
    print("minlen: ", minlen)
    best = ''
    best_val = 0
    if maxlen != 0:
        minlen = maxlen
    for i in range(len(logPath)):
        if 'his' in logPath[i] and 'maxmem' in logPath[i]:
            h_m = warm_start[logPath[i]][minlen-1]
        if best_val < warm_start[logPath[i]][minlen-1] and 'ideal' not in label[logPath[i]]:
            best = label[logPath[i]]
            best_val = warm_start[logPath[i]][minlen-1]
        print(logPath[i], warm_start[logPath[i]][minlen-1])
    print("\nbest warmstart rate: ", best, best_val, "increase: ", 100.0 * (best_val - h_m)/h_m, "%")
    plt.legend(loc='upper left')
    plt.title('total warmstart Rate')
    plt.xlabel('Minute')
    plt.ylabel('warmstart Rate')
    plt.grid(True)
    plt.tight_layout()
    plt.savefig(f"./pkg/result/multi-warmstart-{str(id)}.png")
    plt.close()    

def draw_cost():
    plt.figure(figsize=(15, 25))
    plt.rcParams.update({'font.size': 30})
    plt.subplot(2, 1, 1)
    minlen = 1000000000
    h_m = 1
    for i in range(len(logPath)):
        plt.plot(timeCost[logPath[i]], label=label[logPath[i]], linewidth=2)
        minlen = min(minlen, len(timeCost[logPath[i]]))
    print("minlen: ", minlen)
    best = ''
    best_val = 1000000000000
    if maxlen != 0:
        minlen = maxlen
    for i in range(len(logPath)):
        if 'his' in logPath[i] and 'maxmem' in logPath[i]:
            h_m = timeCost[logPath[i]][minlen-1]
        if best_val > timeCost[logPath[i]][minlen-1] and 'ideal' not in label[logPath[i]]:
            best = label[logPath[i]]
            best_val = timeCost[logPath[i]][minlen-1]
        print(logPath[i], timeCost[logPath[i]][minlen-1], len(timeCost[logPath[i]]))
    print("\nbest timecost: ", best, best_val, "decrease: ", 100.0 * (h_m - best_val)/h_m, "%")
    plt.legend(loc='upper left')
    plt.title('total timecost')
    plt.xlabel('xxx')
    plt.ylabel('timecost Rate')
    
    # plt.subplot(2, 1, 2)
    # minlen = 1000000000
    # for i in range(len(logPath)):
    #     plt.plot(memCost[logPath[i]], label=label[logPath[i]], linewidth=2)
    #     minlen = min(minlen, len(memCost[logPath[i]]))
    # print("minlen: ", minlen)
    # best = ''
    # best_val = 1000000000000
    # h_m = 1
    # if maxlen != 0:
    #     minlen = maxlen
    # for i in range(len(logPath)):
    #     if 'his' in logPath[i] and 'maxmem' in logPath[i]:
    #         h_m = memCost[logPath[i]][minlen-1]
    #     if best_val > memCost[logPath[i]][minlen-1] and 'ideal' not in label[logPath[i]]:
    #         best = label[logPath[i]]
    #         best_val = memCost[logPath[i]][minlen-1]
    #     print(logPath[i], memCost[logPath[i]][minlen-1])
    # print("\nbest memcost: ", best, best_val, "decrease: ", 100.0 * (h_m - best_val)/h_m, "%")
    # plt.legend(loc='upper left')
    # plt.title('total memcost')
    # plt.xlabel('xxx')
    # plt.ylabel('memcost Rate')
    
    # plt.grid(True)
    # plt.tight_layout()
    # plt.savefig(f"./pkg/result/multi-cost-{str(id)}.png")
    # plt.close() 
    
def draw_score1():
    plt.figure(figsize=(20, 30))
    plt.rcParams.update({'font.size': 20})
    for i in range(len(logPath)):
        plt.plot(avg_mem_score[logPath[i]], label=label[logPath[i]])
        plt.plot(avg_time_score[logPath[i]], label=label[logPath[i]])
    plt.legend(loc='upper right')
    plt.title('Average Score')
    plt.xlabel('Minute')
    plt.ylabel('Score (%)')
    plt.savefig(f"./pkg/result/multi-score2-{str(id)}.png")
    plt.close()
    plt.grid(True)

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
        type = path.split('/')[0]
        policy = path.split('/')[1]
        keepAlive = path.split('-')[2]
        memory = path.split('-')[3]
        label[path] = type + '-' + policy + '-' + keepAlive + '-' + memory
    read()
    draw_score()
    draw_cost()
    draw_warmstart()
    print("-----------------------------------------------------------")