# GoAppMation
Portable App Creator in Go

Note: Project is in development so it may change frequently.

This package makes it easy to find latest version of an application found on a website, download the zip to your computer, extract only the necessary files, and then add any files and scripts you need to make the software portable.

This tool will eventually be used to update [golang-portable-windows](https://github.com/josephspurrier/golang-portable-windows) and [surfstack-wamp](https://github.com/josephspurrier/surfstack-wamp) because much of the work to create portable versions of the latest software is manual.

## Example: Build MySQL Portable for Windows

In the config folder, there is a file called: mysql.json. This file contains all the information on how to create a portable distribution of MySQL.

To install, run the following command:
~~~
go get github.com/josephspurrier/goappmation/cmd/goappmation
~~~

Now run goappmation.exe. The portable distribution will be created in a folder called MySQL Portable v5.6.23. The folder will contain a Start.cmd and a Stop.cmd to control MySQL.