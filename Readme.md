### Compile proto buffer
    make
### Deploy on local k8s cluster
    skaffold dev --port-forward
### Introspect server APIs with
    evans --host localhost --port <SERVER_PORT> --reflection repl         
