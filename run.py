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
NUM_TUPLES = 100_000_000
BATCH_SIZE = 5_000

def runSender(numReciever, gomaxprocs, numTuples, batchSize):
  os.system(f"./bin/sender --num={numReciever} --proc={gomaxprocs} -t={numTuples} --batch={batchSize}")

def runReciever(port, gomaxprocs):
  os.system(f"./bin/reciever --port={port} --proc={gomaxprocs}")

def runExperiment(numReciever, gomaxprocs):
  # run receivers
  for i in range(numReciever):
    runReciever(RECIEVER_PORTS[i], gomaxprocs)
  # run sender
  runSender(numReciever, gomaxprocs, NUM_TUPLES, BATCH_SIZE)

def main():
  runExperiment(1, 2)

if __name__=="__main__":
  main()