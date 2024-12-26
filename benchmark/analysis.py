import json
import pandas as pd
import numpy as np
import matplotlib.pyplot as plt

with open("benchmark/bm_results/seq.json", "r") as f:
    seq_result = json.load(f)

with open("benchmark/bm_results/par.json", "r") as f:
    par_result = json.load(f)

with open("benchmark/bm_results/ws.json", "r") as f:
    ws_result = json.load(f)

average_seq = 0
for run in seq_result:
    average_seq += seq_result[run]
average_seq /= len(seq_result)

average_par = {}
for thread in par_result:
    if thread not in average_par:
        average_par[thread] = 0
    average_par[thread] = np.mean(par_result[thread])

average_ws = {}
for thread in ws_result:
    if thread not in average_ws:
        average_ws[thread] = 0
    average_ws[thread] = np.mean(ws_result[thread])

runs = ["Parallel Run", "Work Stealing Run"]

speedups = [{}, {}]
for i, run in enumerate([average_par, average_ws]):
    for thread in run:
        speedups[i][thread] = average_seq / run[thread]

# plot parallel
fig, ax = plt.subplots()
ax.plot(list(speedups[0].keys()), list(speedups[0].values()))
ax.set_xlabel("Number of threads")
ax.set_ylabel("Speedup")
ax.legend()
ax.set_title("Parallel Speedup Performance for Gravitational RK4 System")
plt.savefig("benchmark/bm_results/parallel.png")

# plot work stealing
fig, ax = plt.subplots()
ax.plot(list(speedups[1].keys()), list(speedups[1].values()))
ax.set_xlabel("Number of threads")
ax.set_ylabel("Speedup")
ax.legend()
ax.set_title("Work Stealing Speedup Performance for Gravitational RK4 System")
plt.savefig("benchmark/bm_results/work_stealing.png")


