import subprocess as sp
import pandas as pd
import time as tm
import psutil
import sys


def ParsingArguments(row):
	ClientCommand = "./main client run "
	SimpleDistributions = ["exponential", "poisson"]

	NewClientCommand = ClientCommand + " --profile_output_overwrite=true" + " --iterations=" + str(row['iterations']) + " --concurrent=" + str(row['concurrent']) + " --distribution=" + str(row['distribution'])

	if str(row['distribution']).lower() in SimpleDistributions:
		NewClientCommand = NewClientCommand + " --distribution_params=" + str(row['P1'])
	else:
		NewClientCommand = NewClientCommand + " --distribution_params=" + str(row['P1']) + "," + str(row['P2'])

	NewClientCommand = NewClientCommand + " --profile_output=" + str(row['outputfile']) + " --experiment_description=\"" + str(row['description']) +"\""
	return NewClientCommand

def StartServer():
	ServerCommand = "./main server run &"
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

	if(len(sys.argv) != 2):
		print "Usage: python run_experiments.py experiments.csv"
	else:
                build_proc = sp.Popen("go build main.go")
                build_proc.communicate()

		experiments = pd.read_csv(str(sys.argv[1]))

		for index, row in experiments.iterrows():

			s_handle = StartServer()

			ClientCommand = ParsingArguments(row)

			StartClient(ClientCommand)

			EndServer(s_handle)


if __name__== "__main__":
	main()


