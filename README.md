# pmhb-redis
An implementation which implements apis for working with Redis

# Getting Started
- Set up Redis cluster: https://github.com/carwestsam/redis-cluster-in-docker-osx

- Get pmhb-redis
    ```bash
    go get https://github.com/phongcao0114/pmhb-redis
    ```

# Running
1. Start the redis cluster
    ```bash
    cd ./carwestsam/redis-cluster-in-docker-osx
    bash start.sh
    ```
2. Run pmhb-redis
3. Make API calls

# APIs

##Set:

    {
    	"request_body":{
    		"key":"redis-key-001",
    		"employee":{
    			"name":"john",
    			"position":"director"
    		},
    		"expiry_time":900
    	}
    }

##Get:

        {
    	"request_body":{
    		"key":"redis-key-001"
    	}
    }