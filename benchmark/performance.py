import subprocess
import time
import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import json

def run_seq(num_experiments=5):
    # print("Running sequential experiments")
    res = {}
    for i in range(num_experiments):
        # print(f"Running SEQ experiment {i}")
        result = subprocess.run(
            ["go", "run", "editor/editor.go", "benchmark/in_out/init_values.txt", "seq", "1"],
            capture_output=True, text=True)
        if result.returncode != 0:
            raise RuntimeError("Go program failed:", result.stderr)
        res[i] = float(result.stdout)
    return res

def run_par(threads, num_experiments=5):
    res = {}
    for i in range(num_experiments):
        for thread in threads:
            # print(f"Running {thread} threads for PAR experiment {i}")
            result = subprocess.run(
                ["go", "run", "editor/editor.go", "benchmark/in_out/init_values.txt", "par", str(thread)],
                capture_output=True, text=True)
            if result.returncode != 0:
                raise RuntimeError("Go program failed:", result.stderr)
            if thread not in res:
                res[thread] = []
            res[thread].append(float(result.stdout))
    return res

def run_ws(threads, num_experiments=5):
    res = {}
    for i in range(num_experiments):
        for thread in threads:
            # print(f"Running {thread} threads for WS experiment {i}")
            result = subprocess.run(
                ["go", "run", "editor/editor.go", "benchmark/in_out/init_values.txt", "ws", str(thread)],
                capture_output=True, text=True)
            if result.returncode != 0:
                raise RuntimeError("Go program failed:", result.stderr)
            if thread not in res:
                res[thread] = []
            res[thread].append(float(result.stdout))
    return res

threads = [2, 4, 6, 8, 12]

result = run_seq()
with open("benchmark/bm_results/seq.json", "w") as f:
    json.dump(result, f)

result = run_par(threads)
with open("benchmark/bm_results/par.json", "w") as f:
    json.dump(result, f)

result = run_ws(threads)
with open("benchmark/bm_results/ws.json", "w") as f:
    json.dump(result, f)

