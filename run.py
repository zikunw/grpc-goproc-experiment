import os
import subprocess
import time

EXPERIMENTS = [
  {"numReciever": 1, "gomaxprocs": 1},
  {"numReciever": 1, "gomaxprocs": 2},
  {"numReciever": 1, "gomaxprocs": 4},
  {"numReciever": 1, "gomaxprocs": 8},
  {"numReciever": 2, "gomaxprocs": 1},
  {"numReciever": 2, "gomaxprocs": 2},
  {"numReciever": 2, "gomaxprocs": 4},
  {"numReciever": 2, "gomaxprocs": 8},
  {"numReciever": 4, "gomaxprocs": 1},
  {"numReciever": 4, "gomaxprocs": 2},
  {"numReciever": 4, "gomaxprocs": 4},
  {"numReciever": 4, "gomaxprocs": 8},
  {"numReciever": 8, "gomaxprocs": 1},
  {"numReciever": 8, "gomaxprocs": 2},
  {"numReciever": 8, "gomaxprocs": 4},
  {"numReciever": 8, "gomaxprocs": 8},
]
RECIEVER_PORTS = [
  10000,
  10001,
  10002,
  10003,
  10004,
  10005,
  10006,
  10007,
  10008,
]
NUM_TUPLES = 10_000_000
BATCH_SIZE = 5_000
REPEAT = 3
USE_PROFILE = "false"
USE_CPU = False

def runSender(numReciever, gomaxprocs, numTuples, batchSize, useProfile, cpuProfile):
  os.system(f"./bin/sender --num={numReciever} --proc={gomaxprocs} -t={numTuples} --batch={batchSize} --trace={useProfile} --cpuprofile={cpuProfile}")

def runReciever(port, gomaxprocs):
  run_process_in_background(f"./bin/reciever --port={port} --proc={gomaxprocs}")

def runExperiment(numReciever, gomaxprocs, i):
  # run receivers
  for i in range(numReciever):
    runReciever(RECIEVER_PORTS[i], 2)
  # run sender
  cpuprofile = f"cpu_{numReciever}_{gomaxprocs}_{i}" if USE_CPU else ""
  runSender(numReciever, gomaxprocs, NUM_TUPLES, BATCH_SIZE, USE_PROFILE, cpuprofile)


def run_process_in_background(command):
  try:
    # Run the command in the background and hide the output
    subprocess.Popen(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    print(f"Process '{command}' started in the background.")
  except Exception as e:
    print(f"An error occurred when running {command}: {e}")

def main():
  for e in EXPERIMENTS:
    for i in range(REPEAT):
      runExperiment(e["numReciever"], e["gomaxprocs"], i)
      time.sleep(30)

if __name__=="__main__":
  main()