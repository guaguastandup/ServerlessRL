import matplotlib.pyplot as plt
import numpy as np

# 设置数据
# x_labels = ['200', '400', '600' , '800']
x_labels=['200', '400', '600']
algorithm_names = ['LRU-Fixed', 'MaxMem-Fixed', 'FISH-Fixed', 
                   'Random-HIST', 'LRU-HIST', 'MaxMem-HIST', 'FISH-HIST']

random = np.array([
    [252899.8244],
    [192592.7521],
    [132599.2333],
    # [80853.8078],   
])

hit_rates = np.array([
    [228186.6745, 251155.8966, 227621.9509,
     251118.9973, 230078.0975, 248635.1764, 223934.8236],  # 200G
    [166937.3455, 181781.1119, 157042.8451, 
     190385.3489, 164066.1279, 178846.5897, 155469.9084],  # 400G
    [116563.6632, 103440.8215, 91184.7234, 
     131233.9758, 115076.1663, 102344.7469, 91217.6429],  # 800G
    # [71314.6459, 51527.8406 , 48211.5132, 
    #  80423.25, 70701.7867, 51430.3036, 48034.7527]  # 1000G
])

for i in range(len(hit_rates)):
    for j in range(len(hit_rates[i])):
        hit_rates[i][j] = 100.0 * (random[i] - hit_rates[i][j]) / random[i]
        # hit_rates[i][j] = 100.0 * (random[i] - hit_rates[i][j]) / random[i]
        # print(hit_rates[i][j])
        
n_groups = len(x_labels)
n_algorithms = len(algorithm_names)
index = np.arange(n_groups)
bar_width = 0.12
opacity = 1.0

# 绘制柱状图
fig, ax = plt.subplots(figsize=(15, 8))
hatches = ['/', '\\', '-', '|', '.', '+', 'x', 'X']   # 不同的填充样式
for i in range(n_algorithms):
    if i == n_algorithms - 1:
        bars = plt.bar(index + i * bar_width, hit_rates[:, i], bar_width, alpha=opacity, 
                       label=algorithm_names[i], color='black')  # 实心填充
    else:
        bars = plt.bar(index + i * bar_width, hit_rates[:, i], bar_width, alpha=opacity, 
                   label=algorithm_names[i], color='white', edgecolor='black', hatch=hatches[i % len(hatches)])
    # plt.bar(index + i * bar_width, hit_rates[:, i], bar_width, alpha=opacity, label=algorithm_names[i])

# 添加图例、标题和轴标签
plt.xlabel('Memory size (GB)', fontsize=30)  # 调整轴标签字体大小
plt.ylabel('Relative Saving Time (%)', fontsize=30)  # 调整轴标签字体大小
# plt.title('Cold-Start Time of Algorithms', fontsize=30)

plt.xticks(index + bar_width * (n_algorithms - 1) / 2, x_labels, fontsize=25)
plt.yticks(fontsize=25)
plt.legend(fontsize=25)  # 调整图例字体大小

# 优化布局并显示图表
plt.tight_layout()
plt.savefig("./pkg/result/bar_timecost.pdf")
