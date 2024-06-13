## This project aims to provide functionality to solve order packaging problem.

### *Requirements
`go version 1.22`
because I'm using new functionalities provided by the standard `net/http` library. \
`curl or postman` to be able to call the endpoints.

### Install & run
- download the code locally
- go to `{repo_dir}/cmd/api`
- run `go run .` and you should see the following message `Listening on port 8282` \
Note: you can change this port inside `config\env.go` file or by executing `export CUSTOM_PORT=your_port`
### Exposed APIs
- **addPack**: used to add new packaging size, please specify a list of numeric values inside `--data` param \
replace `{"sizes":[values_here]}` with the value that you want.
```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"sizes":[250,500,1000,2000,5000]}' \
  http://localhost:8282/pack
```
Response: `{"status":"success"}` or `{"error":"some error"}`

- **removePack**: used to remove packaging size \
replace `{size}` with the size that you want to remove, eg. 5000
```
curl --request "DELETE" http://localhost:8282/pack/{size}
```
Response: `{"status":"success"}` or `{"error":"some error"}`

- **removePacks**: used to remove all packaging sizes, becomes handy when you'd want to clear DB.
```
curl --request "DELETE" http://localhost:8282/packs
```
Response: `{"status":"success"}` or `{"error":"some error"}`

- **getOrderPackaging**: used retrieve packaging configuration for given size \
  replace `{size}` with the size that you want to compute configuration, eg. 12001
```
curl --request "GET" http://localhost:8282/order/{size}
```
Response: `{"2000":1,"250":1,"5000":2}` \
And this translate into: \
2 pack(s) of 5000 \
1 pack(s) of 1000 \
1 pack(s) of 250 
