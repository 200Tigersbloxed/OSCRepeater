# OSCRepeater
A tool to repeat OSC messages written in Go

## How to use

Download the latest release, and launch the executable.

> ℹ️ If you're on macOS or Linux, it should still work native without wine, you'll just have to compile builds yourself. ℹ️

Once you open the application, you'll see a view like below

![image](https://user-images.githubusercontent.com/45884377/192105715-bba6f57d-a20e-4962-9554-44a3e7916f32.png)

Here's an explanation of each category in the program,

### Applications

All of the applications to listen for repeating. Clicking a button will open an existing Application.

![image](https://user-images.githubusercontent.com/45884377/192080246-0b06acd7-4881-4ca4-999b-4beae55869f0.png)

+ **ApplicationName**
  + The name of your Applications
+ **ListenAddress**
  + The Address of where your messages are going to be forwarded to
+ **ListenPort**
  + The Port to where the messages will be forwarded to
+ **SendHost**
  + Where to listen for incoming messages from the application
+ **SAVE**
  + Saves the current data to config
+ **DELETE**
  + Deletes the application from config

___

### Actions

A list of all Actions that can be done

+ **Start**
  + Starts all OSC components
+ **Stop**
  + Stops all OSC components
+ **Create Application**
  + Creates a new Application
+ **Reload Config**
  + Reloads from the `config.json` file

___

+ **Status**
  + Lists the current OSC status
