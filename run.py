import os
import subprocess

EXPERIMENTS = [
  {
    "numReciever": 1,
    "gomaxprocs": 1,
  }
]
RECIEVER_PORTS = [
  10000,
  10001,
  10002,
  10003,
  10004,
  10005
]
NUM_TUPLES = 1_000_000
BATCH_SIZE = 5_000

def runSender(numReciever, gomaxprocs, numTuples, batchSize):
  os.system(f"./bin/sender --num={numReciever} --proc={gomaxprocs} -t={numTuples} --batch={batchSize}")

def runReciever(port, gomaxprocs):
  run_process_in_background(f"./bin/reciever --port={port} --proc={gomaxprocs}")

def runExperiment(numReciever, gomaxprocs):
  # run receivers
  for i in range(numReciever):
    runReciever(RECIEVER_PORTS[i], 2)
  # run sender
  runSender(numReciever, gomaxprocs, NUM_TUPLES, BATCH_SIZE)

def run_process_in_background(command):
  """
  Run a process in the background.
  """

  try:
    # Run the command in the background and hide the output
    subprocess.Popen(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    print(f"Process '{command}' started in the background.")
  except Exception as e:
    print(f"An error occurred when running {command}: {e}")

def main():
  runExperiment(1, 2)

if __name__=="__main__":
  main()