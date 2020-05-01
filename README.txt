#Authors: Zach Luciano, Tanner Halcumb, Nick Mladineo
#Date: April 30, 2020
#Project: GoScraper

#Environment
  Machine must be a Windows machine with Go installed and configured, see https://golang.org/doc/install
  for install and configuration information. Testing did not prove successful on Windows machines without
  Go installed and other OS systems as the executable is ".exe".

#Running
  With all source files downloaded, go to the "src/" folder. The following command will need to be run
  to gather the necessary packages:
===
$ go get github.com/asticode/go-astikit github.com/asticode/go-astilectron github.com/bbalet/stopwords github.com/gocolly/colly github.com/gorilla/mux github.com/lib/pq
===

  Once the packages are installed, the build can be done for the executable.
===
$ go build
===

  The build command will produce "src.exe" simply run this file to execute the program.