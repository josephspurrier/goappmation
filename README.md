# GoAppMation
Portable App Creator in Go

Note: Project is in development so it may change frequently.

This package makes it easy to find latest version of an application found on a website, download the zip to your computer, extract only the necessary files, and then add any files and scripts you need to make the software portable.

This tool is currently used to update [golang-portable-windows](https://github.com/josephspurrier/golang-portable-windows) and will be eventually used to update [surfstack-wamp](https://github.com/josephspurrier/surfstack-wamp) because much of the work to create portable versions of the latest software is manual.

## Example: Build MySQL Portable for Windows

In the **project/MySQL Portable** folder, there are two json config files. These files contains all the instructions on how to create a portable distribution of MySQL.

To install, run the following command:
~~~
go get github.com/josephspurrier/goappmation/cmd/goappmation
~~~

Then, run one of these commands depending on how you want to specify the version:

~~~
# MySQL v5.7.12 - version specified inside the config file
goappmation.exe "../../project/MySQL Portable/MySQL Portable v5.7.12.json"

# MySQL v5.7.12 - version specified outside the config file
goappmation.exe "../../project/MySQL Portable/MySQL Portable v5.7.X.json" -version=5.7.12
~~~

The portable distribution will be created in a folder called: MySQL Portable v5.7.12

The folder will contain the following batch scripts:

* _Initialize.cmd - Run first to create the data directory before starting MySQL
* _Start.cmd      - Run to start MySQL
* _Stop.cmd       - Run to stop MySQL

## Variables in config

You can use the **{{VERSION}}** variable in **ApplicationName** and **DownloadUrl**.

## Goals of this Project

### Completed

* Download a ZIP file from a URL using a version number
* Extract certain files from a URL using regular expressions
* Remove the root folder from a ZIP archive
* Create files to do certain things like run the application
* Get the version number of the latest version
* Move/rename files and folders
* Pass version that overwrites version in the config file

### Next Goals

* Add ability to download files that are not zipped
* Add ability to edit files
* Add ability to append to files
* Handle migrations between different versions
* Do diffs between text files for configs
* Maybe create a simple script language