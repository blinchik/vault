

```bash
go run main.go -enable <vault address> 8200 <vault token> pki pki 72h 72h
go run main.go -enable <vault address> 8200 <vault token> pki pki_int 72h 72h

go run main.go -gr <vault address> 8200 <vault token> blinchik.consul

go run main.go -setURL <vault address> 8200 <vault token> http://<vault address>:8200/v1/pki/crl http://<vault address>:8200/v1/pki/ca


go run main.go -gi <vault address> 8200 <vault token> "blinchik.consul Intermediate Authority" pki_int


after this we create a role 

go run main.go -cr <vault address> 8200 <vault token> pki_int tlsAdmin true "blinchik.consul"
```
