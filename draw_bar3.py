import matplotlib.pyplot as plt
import numpy as np

# 设置数据
# x_labels = ['200G', '400G', '600G' , '800G']
x_labels = ['200', '400', '600' , '800']
algorithm_names = ['LRU-Fixed', 'MaxMem-Fixed', 'FISH-Fixed', 
                   'Random-HIST', 'LRU-HIST', 'MaxMem-HIST', 'FISH-HIST']

random = np.array([
    [0.973],
    [0.884],
    [0.7836],
    [0.6767],   
])

hit_rates = np.array([
    [0.9657, 0.9591, 0.9393, 
     0.9721, 0.9681, 0.9564, 0.9367],  # 200G
    [0.8726, 0.8492, 0.8335, 
     0.881, 0.8705, 0.8458, 0.8316],  # 400G
    [0.7718, 0.7216, 0.7146, 
     0.7822, 0.771, 0.7204, 0.7146],  # 800G
    [0.6889, 0.5943, 0.5941, 
    0.676, 0.6887, 0.5873, 0.5873]  # 1000G
])

for i in range(len(hit_rates)):
    for j in range(len(hit_rates[i])):
        hit_rates[i][j] = 100.0 * (random[i] - hit_rates[i][j]) / random[i]
        
n_groups = len(x_labels)
n_algorithms = len(algorithm_names)
index = np.arange(n_groups)
bar_width = 0.11
opacity = 0.9

# 绘制柱状图
fig, ax = plt.subplots()
for i in range(n_algorithms):
    plt.bar(index + i * bar_width, hit_rates[:, i], bar_width, alpha=opacity, label=algorithm_names[i])

# 添加图例、标题和轴标签
plt.xlabel('Memory size(GB)')
plt.ylabel('Relative cold start rate(%)')
plt.title('Relative cold start rate of different algorithms')
plt.xticks(index + bar_width * (n_algorithms - 1) / 2, x_labels)
plt.legend()

# 优化布局并显示图表
plt.tight_layout()
# plt.show()
plt.savefig("./pkg/result/bar_cold_satrt.pdf")
