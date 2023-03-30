## POC to make DAGs for decision routing
This is a POC on how to create DAGs and use that to drive 
routing decisions

### How to build
You can use docker to build and test this locally.
See the Makefile for rules 

#### building and running
Use the following commands to build and run docker
```
make docker-build && make docker-run
```

Use Ctrl+C to end the container once testing is done

Use this to clean up post-testing
```
make docker-clean
```

### Sample decision tree
![alt text for screen readers](tree.png "sample tree used in main.go")

The red arrows are Choice A

The blue is Choice B

and green is Choice C

### Endpoints
#### POST /next
 ```
 curl --location 'localhost:8080/next' \
--header 'Content-Type: application/json' \
--data '{
    "current_node": "specialty_node",
    "input_kvp": {
        "choice": "a"
    }
}'
```

which takes in 2 values:
- current node id
- key-value pair with whatever inputs we wanna share with the server

the decision of which node to go to would be made using these 2 inputs

Sample Response
```
next node id: services_node
```

#### GET /node/{id}
This endpoint is used to fetch metadata about a specific node.
Can be invoked using
```
curl --location 'localhost:8080/node/specialty_node'
```
 which returns the response:
 ```
route: /specialties
Method: GET
 ```

 ### Debugging Issues
 #### How to run locally without Docker
 ```
 make run
 ```

 #### How to ensure 8080 port is free
See https://stackoverflow.com/a/36876427/224640