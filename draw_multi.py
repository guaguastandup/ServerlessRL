import pandas as pd
import matplotlib.pyplot as plt
import sys
import numpy as np
mem_running_usage, mem_occupy_usage = {}, {}
mem_score, time_score = {}, {}
warm_start_rate, cdf_warmstart = {}, {}
app_mem_score, app_time_score = {}, {}
avg_coldstart, avg_mem_score, avg_time_score = {}, {}, {}

logPath = []

keepAliveList=[5]
policyList=[
    'maxmem',    # 第1
    # 'maxmem2',   # 第1
    # 'score',      # 第2
    # 'maxUsage',     # 第3
    # 'maxColdStart', # 第3
    # 'score3', # 第5
    # 'random',       # 第7
    # 'minColdStart', # 第8
    # 'lru',       # 第9
    # 'lfu',
    # 'mru',
]
policyList2=[
    'maxmem',    # 第1
    # 'maxmem2',   # 第1
    # 'score',      # 第2
    # 'maxUsage',     # 第3
    # 'maxColdStart', # 第3
    'score1', # 第4
    'score2', # 第4
    # 'score3', # 第5
    'score4',
    'score5',
    'score6',
    # 'random',       # 第7
    # 'minColdStart', # 第8
    # 'lru',       # 第9
    # 'lfu',
    # 'mru',
]

memoryList=[1000]
arrivalCnt=1

id = 0
if len(sys.argv) > 1:
    id = int(sys.argv[1])
    
label = {}

for k in keepAliveList[0:1]:
    for p in policyList:
        for m in memoryList[0:1]:
            logPath.append('fixed/' + p + '/fixed-' + p + '-' + str(k) + '-' +str(m) + '-' + str(arrivalCnt))

for k in keepAliveList[0:1]:
    for p in policyList2:
        for m in memoryList[0:1]:
            logPath.append('histogram/' + p + '/histogram-' + p + '-' + str(k) + '-' + str(m) + '-' + str(arrivalCnt))
                        
for path in logPath:
    type = path.split('/')[0]
    policy = path.split('/')[1]
    keepAlive = path.split('-')[2]
    memory = path.split('-')[3]
    label[path] = type + '-' + policy + '-' + keepAlive + '-' + memory

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
            avg_coldstart[logPath[i]] = []
            avg_mem_score[logPath[i]] = []
            avg_time_score[logPath[i]] = []
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
    
    # plt.rcParams.update({'font.size': 15})
    # plt.figure(figsize=(40, 20))
    # plt.subplot(3, 2, 1)
    # plt.plot(mem_running_usage[logPath[0]], label='Mem Running')
    # for i in range(len(logPath)):
    #     plt.plot(mem_occupy_usage[logPath[i]], label=label[logPath[i]])
    # plt.legend(loc='upper right')
    # plt.title('MEMOccupyUsage')
    # plt.xlabel('Minute')
    # plt.ylabel('Usage (GB)')
    # plt.grid(True)

    # plt.subplot(3, 2, 2)
    # for i in range(len(logPath)):
    #     for j in range(len(cdf_warmstart[logPath[i]])):
    #         cdf_warmstart[logPath[i]][j] = (1.0 - cdf_warmstart[logPath[i]][j]) * 100
    #     data_sorted = np.sort(cdf_warmstart[logPath[i]])
    #     cdf = np.arange(1, len(data_sorted) + 1) / len(data_sorted)
    #     plt.plot(data_sorted, cdf, label=label[logPath[i]])
    # plt.legend(loc='upper right')
    # plt.title('CDF of ColdStart Rate')
    # plt.xlabel('ColdStart Rate (%)')
    # plt.ylabel('CDF')
    # plt.grid(True)

    # plt.subplot(3, 2, 3)
    # for i in range(len(logPath)):
    #     plt.plot(mem_score[logPath[i]], label=label[logPath[i]])
    # plt.legend(loc='upper right')
    # plt.title('Mem Score')
    # plt.xlabel('Minute')
    # plt.ylabel('Score (%)')
    # plt.grid(True)

    # plt.subplot(3, 2, 4)
    # for i in range(len(logPath)):
    #     plt.plot(time_score[logPath[i]], label=label[logPath[i]])
    # plt.legend(loc='upper right')
    # plt.title('Time Score')
    # plt.xlabel('Minute')
    # plt.ylabel('Score (%)')
    # plt.grid(True)


    # plt.subplot(3, 2, 5)
    # for i in range(len(logPath)):
    #     data_sorted = np.sort(app_mem_score[logPath[i]])
    #     cdf = np.arange(1, len(data_sorted) + 1) / len(data_sorted)
    #     plt.plot(data_sorted, cdf, label=label[logPath[i]])
    # plt.legend(loc='upper right')
    # plt.title('CDF of Mem Score')
    # plt.xlabel('Mem Score (%)')
    # plt.ylabel('CDF')
    
    # plt.subplot(3, 2, 6)
    # for i in range(len(logPath)):
    #     data_sorted = np.sort(app_time_score[logPath[i]])
    #     cdf = np.arange(1, len(data_sorted) + 1) / len(data_sorted)
    #     plt.plot(data_sorted, cdf, label=label[logPath[i]])
    # plt.legend(loc='upper right')
    # plt.title('CDF of Time Score')
    # plt.xlabel('Time Score (%)')
    # plt.ylabel('CDF')
    
    # plt.tight_layout()
    # plt.savefig(f"./pkg/result/multi-{str(id)}.png")
    # plt.close()

    plt.figure(figsize=(20, 30))
    # plt.subplot(1, 2, 1)
    # for i in range(len(logPath)):
    #     plt.plot(avg_mem_score[logPath[i]], label=label[logPath[i]])
    #     plt.plot(avg_time_score[logPath[i]], label=label[logPath[i]])
    # plt.legend(loc='upper right')
    # plt.title('Average Score')
    # plt.xlabel('Minute')
    # plt.ylabel('Score (%)')
    # plt.grid(True)
    
    # plt.subplot(1, 2, 2)
    plt.rcParams.update({'font.size': 20})
    for i in range(len(logPath)):
        plt.plot(avg_coldstart[logPath[i]], label=label[logPath[i]])
    plt.legend(loc='upper left')
    plt.title('Average ColdStart Rate')
    plt.xlabel('Minute')
    plt.ylabel('ColdStart Rate')
    plt.grid(True)
    
    plt.tight_layout()
    plt.savefig(f"./pkg/result/multi-score-{str(id)}.png")
    plt.close()

    # plt.figure(figsize=(15, 25))
    # for i in range(len(logPath)):
    #     data_sorted = np.sort(cdf_warmstart[logPath[i]])
    #     cdf = np.arange(1, len(data_sorted) + 1) / len(data_sorted)
    #     plt.plot(data_sorted, cdf, label=label[logPath[i]])
    # plt.legend(loc='upper right')
    # plt.title('CDF of ColdStart Rate')
    # plt.xlabel('ColdStart Rate (%)')
    # plt.ylabel('CDF')
    # plt.grid(True)
    # plt.savefig(f"./pkg/result/cdf_multi-{str(id)}.png")
    # plt.close()

    # fig, ax = plt.subplots()
    # fig.set_size_inches(30, 20)
    # policyArray = {}
    # policies = []
    # for file in logPath:
    #     dot10, dot25, dot50, dot75, dot90, dot100 = 0, 0, 0, 0, 0, 0
    #     for i in range(len(cdf_warmstart[file])):
    #         if 0 <= cdf_warmstart[file][i] < 10:
    #             dot10 += 1
    #         elif 10 <= cdf_warmstart[file][i] < 25:
    #             dot25 += 1
    #         elif 25 <= cdf_warmstart[file][i] < 50:
    #             dot50 += 1
    #         elif 50 <= cdf_warmstart[file][i] < 75:
    #             dot75 += 1
    #         elif 75 <= cdf_warmstart[file][i] < 90:
    #             dot90 += 1
    #         elif 90 <= cdf_warmstart[file][i] <= 100:
    #             dot100 += 1
    #     avg = sum(cdf_warmstart[file]) / len(cdf_warmstart[file])
    #     # print(label[file], avg)
    #     # print(label[file], dot10, dot25, dot50, dot75, dot90, dot100, '\n')
    #     tag = file.split('-')[1]
    #     tag = tag.split('/')[0]
    #     if tag not in policyArray:
    #         policyArray[tag] = []
    #         policies.append(tag)
    #     policyArray[tag].append(avg)
    #     print(file.split('/')[0], file.split('/')[1], avg, '\n')
    
    # width = 0.35
    # x = np.arange(len(policies))

    # # 为每个策略绘制并列的柱状图
    # for i, policy in enumerate(policies):
    #     values = policyArray[policy]
    #     ax.bar(x[i] - width/2, values[0], width)
    #     ax.bar(x[i] + width/2, values[1], width)
        
    # ax.set_ylabel('ColdStart Rate (%)')
    # ax.set_title('ColdStart Rate by Policy')
    # ax.set_xticks(x)
    # ax.set_xticklabels(policies)
    # ax.legend()
    # plt.grid(True)
    # plt.savefig(f"./pkg/result/bar_multi-{str(id)}.png")
    # plt.close()
draw()