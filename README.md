# Node LoRaWAN Parameter configurator

### Note
- The program waits for the keyword `>CONF` from the device to start the configuration process, each keyword ends with `\r\n`.

| Keyword | Description|
|---------|------------|
|`>CONF\r\n`| start process|
|`>DEVEUI\r\n`| configure device address 16hex|
|`>APPEUI\r\n`| configure app eui 16hex|
|`>APPKEY\r\n`| configure app key 32hex|
|`>FINISH\r\n`| finish process|

### Build
You must have "gcc" installed before you go build `go build`.

### Author 
Diego Hinojosa Cordova