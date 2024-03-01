import matplotlib.pyplot as plt
import numpy as np

# 设置数据
# x_labels = ['200', '400', '600' , '800']
x_labels = ['200', '400', '600']
algorithm_names = ['Average Cold-Start Decrease', 'Cold-Start Time Decrease', 'Overall Warm-Start Increase']

random = np.array([
    [0.9564, 248635.1764, 11.8537],
    [0.8458, 178846.5897, 46.2776],
    [0.7216, 103440.8215, 71.6472],
    # [0.5943, 51527.8406, 86.9468],   
])

hit_rates = np.array([
    [0.9389, 226290.3766, 16.7054],
    [0.8329, 156784.8216, 49.4621],
    [0.7171, 92227.6312, 72.2805,],
    # [0.5944, 49159.7845, 86.264],
])
for i in range(len(hit_rates)):
    for j in range(len(hit_rates[i])):
        hit_rates[i][j] = 100.0 * (random[i][j] - hit_rates[i][j]) / random[i][j]
        if hit_rates[i][j] < 0:
            hit_rates[i][j] = -hit_rates[i][j]
n_groups = len(x_labels)
n_algorithms = len(algorithm_names)
index = np.arange(n_groups)
bar_width = 0.18
opacity = 1.0
# 绘制柱状图
fig, ax = plt.subplots(figsize=(10, 6))
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
plt.ylabel('Relative Value (%)', fontsize=30)  # 调整轴标签字体大小
# plt.title('Cold-Start Time of Algorithms', fontsize=30)
plt.xticks(index + bar_width * (n_algorithms - 1) / 2, x_labels, fontsize=25)
plt.yticks(fontsize=25)
plt.legend(fontsize=25)  # 调整图例字体大小
# 优化布局并显示图表
plt.tight_layout()
plt.savefig("./pkg/result/bar_memory.pdf")
