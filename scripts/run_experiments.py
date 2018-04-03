import subprocess as sp
import csv
import time as tm
import psutil
import sys


def ParsingArguments(row):
	ClientCommand = "./main client run "
	SimpleDistributions = ["exponential", "poisson"]

	NewClientCommand = ClientCommand + " --profile_output_overwrite=true" + " --iterations=" + str(row['iterations']) + " --concurrent=" + str(row['concurrent']) + " --distribution=" + str(row['distribution'])
        if row['model_percentage'] is not None:
            NewClientCommand = NewClientCommand + " --model_percentage=" + str(row['model_percentage'])
	if str(row['distribution']).lower() in SimpleDistributions:
		NewClientCommand = NewClientCommand + " --distribution_params=" + str(row['P1'])
	else:
		NewClientCommand = NewClientCommand + " --distribution_params=" + str(row['P1']) + "," + str(row['P2'])

	NewClientCommand = NewClientCommand + " --profile_output=" + str(row['outputfile']) + " --experiment_description='" + str(row['description']) +"'"
	return NewClientCommand

def StartServer(eviction_policy):
	ServerCommand = "./main server run --eviction=" + eviction_policy + " &"
	server_process = sp.Popen(ServerCommand.split())
	print "Starting server"
	tm.sleep(10)
	return server_process

def EndServer(server_process):
	Kill = "kill -9 "
	PROCNAME = "uprd"
	print "Ending Server"
	server_process.kill()
	for proc in psutil.process_iter():
    		if proc.name() == PROCNAME:
       			Kill = Kill + str(proc.pid)
        		break

	sp.Popen(Kill.split())

def StartClient(ClientCommand):
	print "Starting Client:"
	print ClientCommand
	client_process = sp.Popen(ClientCommand.split())
	client_process.communicate()


def main():

	if len(sys.argv) != 2 and len(sys.argv) != 3:
		print "Usage: python run_experiments.py experiments.csv [server_policy]"
		sys.exit()

	server_policy = "lru"
	if len(sys.argv) == 3:
		server_policy = str(sys.argv[2])

	build_proc = sp.Popen("go build main.go".split())
	build_proc.communicate()

	with open(str(sys.argv[1]), 'r') as csvfile:
		experiments = csv.DictReader(csvfile)

		for row in experiments:
			if row["outputfile"].startswith("//"):
				continue
			s_handle = StartServer(server_policy)

			ClientCommand = ParsingArguments(row)

			StartClient(ClientCommand)

			EndServer(s_handle)


if __name__== "__main__":
	main()


