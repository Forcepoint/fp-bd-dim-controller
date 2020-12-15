# Controller Installation

Set up machine, install Docker and Docker-Compose

* Docker login:
    * docker.frcpnt.com
        * user.name@forcepoint.com
        * docker.frcpnt.com password

* Copy the docker-compose file from each projects repository into a folder
    * For the controller, the only thing that *MAY* be configured are the passwords, but they have defaults and don't need to be.
    * So in the case that nothing needs to be changed, you can just run the command `docker-compose up -d` in the folder that contains the `docker-compose.yml` file. This should start the database, controller and UI modules.
    * When the controller and UI are up and running, you can now begin to add sub-modules.
    
    * Caveats
        * controller requires folder in its root called `config` with a file called `config.yml`. This file should contain the following:
    ```
      internaltoken: <generated on startup>
      jwt-secret-key: super_secret_key
      local-port: :8080
      logfile: logging
      loglevel: trace
    ```
  
    * Ports and Deployment to testing/production
        * In the `docker-compose` file, the port mappings are 1:1 to the host port to aid in testing, e.g. 8080 -> 8080, in production this can be changed to just expose the ports for the controller and DB internally and only map the UI port to the host port.
        Example of testing configuration:
        ```
      ...
      volumes:
            - ${PWD}/config:/config
      ports:
            - 8080:8080
      networks:
      ...
      ``` 
      Example of production configuration:
      ```
      ...
      volumes:
          - ${PWD}/config:/config
      expose:
          - "8080"
      networks:
      ...
    
* Adding sub-modules
    * Get the docker-compose file for the sub-module project, save it in a folder named after the project (not essential, but easier to find what you are looking for).
    * In the `docker-compose.yml` for each module set `HOST_DOMAIN` to "localhost" and for the `INTERNAL_TOKEN` variable, go to `localhost:8081`, log in and go to the settings icon in the top right, select "Generate Key" and copy that key into the field for `INTERNAL_TOKEN`
    * Next run `docker-compose up -d` and then you should be able to see the new module reflected in the UI after you refresh the page.
    * If you don't see it in the UI, you can run `docker ps` in the command line and check if it is there, if you see it there is probably a problem with the `INTERNAL_TOKEN`, make sure it is present and copied correctly in the `docker-compose` file.
        * Run `docker-compose down` and `docker-compose up -d` and it should be working.
        
    * Caveats
        * The FP-Import module requires a folder called `lists` inits root containing files called `urls.txt` and `blocklist.txt`
        * The AWS Guard Duty module requires a folder called `config` in its root containing a file called `config.yml` with the following:
        ```
      url-token: <Generated at startup>
      ```