# Scatter-Gather

This project is a demonstration of the scatter-gather pattern, commonly used in big data processing to efficiently handle large datasets.

It consists of a distributed processing system where tasks are scattered to multiple worker nodes, processed in parallel, and then gathered back to be aggregated by the orchestrator.

In this project, the orchestrator expects word queries, it then splits the words according to the amount of workers and sends them to them. The workers then read all text files in the `public` diretory and return to the orchestrator the amount of times a specific word appeared per file.

## Requirements

- Go 1.19 or above

## Running the project


If you wish to change the orchestrator proccess port or the amount and/or ports of the worker processes, just change the following variables in the `run.sh` script.
```bash
ORCHESTRATOR_PORT=8080
WORKER_PORTS=(8081 8082 8083 8084)
```

To start the system, run the `run.sh` bash script.
```bash
sh run.sh
```

To connect and query the orchestrator, assuming it is running on port 8080, you can use the `telnet` command.
```bash
telnet localhost 8080
```

To kill all running processes, in the root of the project run.
```bash
make stop
```
