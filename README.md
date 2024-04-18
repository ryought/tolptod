# Tolptod

interactive dotplot


## Usage

Download binary from https://github.com/ryought/tolptod/releases/

```
$ ssh myServer -L 8080:localhost:8080  # if you want to run server on remote machine
$ ./tolptod x.fa y.fa
Parsing x.fa
Parsing y.fa
Building suffix array...
Done
Server running on :8080...
```

Then access to localhost:8080 on your browser

![](/docs/usage.png)

## Develop

frontend
React+typescript
Request to server (with resolution and region)
Draw canvas by match points

backend
Golang
Load fasta and create suffix array.
Query sequences to find matches.
