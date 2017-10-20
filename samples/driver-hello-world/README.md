# driver-hello-world

A hello world driver for go


# To use 

Clone the git repo into the databox root directory 

     git clone https://github.com/me-box/lib-go-databox.git

Run the following commands in a terminal 

```
cd ./lib-go-databox/samples/
docker build -t driver-hello-world .
```
    
 Then run (If databox is not running)
 
      ./databox-start 
      

Finaly upload the manifest file:
 
     go to http://127.0.0.1:8181 in a web browser
     
     select upload and choose /lib-go-databox/samples/driver-hello-world/databox-manifest.json
     
     
 driver-hello-world sould than be available to install in the app store 
 
   

