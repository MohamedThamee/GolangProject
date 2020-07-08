# GolangProject
This is a project developed for Harman interview

## About the project

This project has two types of sub projects
1. Rest API execution
2. Command Line execution

### Rest API execution

This project includes Encode.go file.
Once we run this go file, it starts the server on localhost:5000
We can get the required response through this server

### Command Line execution
This project includes CmdLine Project/cmdProject.go.
Once we run this go file, it starts to run and we can interact with command line.
    

 # Steps to run server for Rest API

  - Run the Encode.go file using the command **go run Encode.go**
  - You can access the rest APIs using the url **http://localhost:5000/**, Once the server is started. 
  - To get the nearest three available Petrol station, Restaurant and Shopping details, use the path `api/{Location}`
        Example:
            - http://localhost:5000/api/sunnyvale
  
  # Steps to run in Commmand line

  - Run the cmdProject.go file which is inside the folder CmdLine Project using the command **go run cmdProject.go**
  -  Once it is started, you will be asked to enter your location and press enter
  -  Now you will be getting the nearest three availabe places for Petrol station, Restaurant and Shopping

  # Included My Resume
  I have included my Resume with this project for your easy reference