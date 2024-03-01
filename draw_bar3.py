import matplotlib.pyplot as plt
import numpy as np

# 设置数据
# x_labels = ['200', '400', '600' , '800']
x_labels=['200', '400', '600']
algorithm_names = ['LRU-Fixed', 'MaxMem-Fixed', 'FISH-Fixed', 
                   'Random-HIST', 'LRU-HIST', 'MaxMem-HIST', 'FISH-HIST']

random = np.array([
    [7.7271],
    [36.722],
    [61.7141],
    # [79.8891],   
])

hit_rates = np.array([
    [8.1258, 11.0253, 16.5055, 8.0239, 7.4587, 11.8537, 17.5185],  # 200G
    [35.6962, 45.5777, 49.4143, 37.3216, 35.8957, 46.2776, 49.7127],
    [64.3574, 71.6472, 72.1405, 62.0337, 64.5247, 71.8448, 72.0757],
    # [81.2925, 86.9468, 86.6332, 79.9402, 81.2521, 86.6394, 86.2488],
])

# 计算相对冷启动率
relative_cold_start_rates = np.zeros_like(hit_rates)
for i in range(len(hit_rates)):
    for j in range(len(hit_rates[i])):
        rate = 100.0 * (hit_rates[i][j] - random[i]) / random[i]
        relative_cold_start_rates[i][j] = max(rate, 0)  # 确保不显示负值的柱子

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
        bars = plt.bar(index + i * bar_width, relative_cold_start_rates[:, i], bar_width, alpha=opacity, 
                       label=algorithm_names[i], color='black')  # 实心填充
    else:
        bars = plt.bar(index + i * bar_width, relative_cold_start_rates[:, i], bar_width, alpha=opacity, 
                   label=algorithm_names[i], color='white', edgecolor='black', hatch=hatches[i % len(hatches)])

    # 对于每个柱子，如果原始数据是负值，则在柱子上方显示该负值
    for rect, original_rate in zip(bars, 100.0 * (hit_rates[:, i] - random.flatten()) / random.flatten()):
        if original_rate < 0:
            height = rect.get_height()
            ax.annotate(f'{original_rate:.1f}',
                        xy=(rect.get_x() + rect.get_width() / 2, height),
                        xytext=(0, 3),  # 3 points vertical offset
                        textcoords="offset points",
                        ha='center', va='bottom',
                        fontsize=15)  # 调整标注字体大小

# 添加图例、标题和轴标签
plt.xlabel('Memory size (GB)', fontsize=30)  # 调整轴标签字体大小
plt.ylabel('Relative Warm-Start Increase (%)', fontsize=30)  # 调整轴标签字体大小
plt.xticks(index + bar_width * (n_algorithms - 1) / 2, x_labels, fontsize=25)  # 调整x轴刻度标签的字体大小
plt.yticks(fontsize=25)  # 调整y轴刻度标签的字体大小
plt.legend(fontsize=25)  # 调整图例字体大小

# 优化布局并显示图表
plt.tight_layout()
plt.savefig("./pkg/result/bar_warmstart.pdf")