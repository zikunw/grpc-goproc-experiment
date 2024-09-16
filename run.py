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

def runSender(numReciever, gomaxprocs, numTuples, batchSize):
  os.system(f"./bin/sender --num={numReciever} --proc={gomaxprocs} -t={numTuples} --batch={batchSize}")

def runReciever(port, gomaxprocs):
  os.system(f"./bin/reciever --port={port} --proc={gomaxprocs}")

def main():


if __name__=="__main__":
  main()