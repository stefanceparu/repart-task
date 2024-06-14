## This project aims to provide functionality to solve order packaging problem.

### *Requirements
`go version 1.22`
because I'm using new functionalities provided by the standard `net/http` library. \
`curl or postman` to be able to call the endpoints.

### Install & run (classic way)
- download the code locally `git clone git@github.com:stefanceparu/repart-task.git`
- in terminal go to `{repo_dir}/cmd/api` & run `go run .`
- or you can use the make commands: `make build` and `make run` inside the root folder.
- you should see the following message `Listening on port 8282` \
Note: you can change this port inside `config\env.go` file or by executing `export CUSTOM_PORT=your_port` then start again the server

### Install & run (using docker)
- download the code locally `git clone git@github.com:stefanceparu/repart-task.git`
- go inside downloaded repo `cd repart-task`
- run `make build-docker`, this will build the docker image
- run `make run-docker`, this will run the docker container using the image created.
- last we need to retrieve the container IP so that we can start making API calls: \
  run `make get-ip` and use this IP address to interact with server. \
  **Note:** If you're using examples below, please make sure that you've replaced the localhost URL with the new address obtained.
- if you want to stop the container use: `make stop-docker`
- if you want to clear container and image use: `make clear-docker`

### Exposed APIs
- **AddPacks [POST /pack]**: used to add new packaging sizes \
  replace `{"sizes":[values_here]}` with the value that you want.
  ```
  curl --header "Content-Type: application/json" \
    --request POST \
    --data '{"sizes":[250,500,1000,2000,5000]}' \
    http://localhost:8282/pack
  ```
  Response: `{"status":"success"}` or `{"error":"some error"}` \


- **RemovePack [DELETE /pack/{size}]**: used to remove packaging size \
replace `{size}` with the size that you want to remove, eg. 5000
  ```
  curl --request "DELETE" http://localhost:8282/pack/{size}
  ```
  Response: `{"status":"success"}` or `{"error":"some error"}` 


- **RemovePacks [DELETE /packs]**: used to remove all packaging sizes, becomes handy when you'd want to clear DB.
  ```
  curl --request "DELETE" http://localhost:8282/packs
  ```
  Response: `{"status":"success"}` or `{"error":"some error"}` 


- **GetOrderPackaging [GET /order/{size}]**: used retrieve packaging configuration for given size \
   replace `{size}` with the size that you want to compute configuration, eg. 12001
  ```
  curl --request "GET" http://localhost:8282/order/{size}
  ```
  Response: `{"2000":1,"250":1,"5000":2}` \
And this translates into: \
2 pack(s) of 5000 \
1 pack(s) of 1000 \
1 pack(s) of 250 
