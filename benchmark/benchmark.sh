#!/bin/bash
#
#SBATCH --mail-user=minhkhoimac@cs.uchicago.edu
#SBATCH --mail-type=ALL
#SBATCH --job-name=proj2_benchmark 
#SBATCH --output=./slurm/out/%j.%N.stdout
#SBATCH --error=./slurm/out/%j.%N.stderr
#SBATCH --chdir=/home/minhkhoimac/parallel-programming/project-2-minhkhoimac/proj2
#SBATCH --partition=debug 
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=16
#SBATCH --mem-per-cpu=900
#SBATCH --exclusive
#SBATCH --time=90:00


module load golang/1.19

python3 benchmark/performance.py

# Check if py1.py ran successfully
if [ $? -eq 0 ]; then
    # Run py2.py if py1.py completed successfully
    python3 benchmark/analysis.py
else
    echo "performance.py encountered an error, so analysis.py will not run."
    exit 1
fi
