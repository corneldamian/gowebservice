# gowebservice

###Simple go web startup project###

If you need to start a rapid web project using go, you can start from here

- logger available for each request in context
- access logger
- in memory session manager
- panic catcher middleware
- config load
- windows/linux/mac os service (you can install it as a service)

check the example dir, it's an working application using this service

is using:
[httpway](https://github.com/corneldamian/httpway.git)
[httpwaymid](https://github.com/corneldamian/httpwaymid.git)
[golog](https://github.com/corneldamian/golog.git)
[service](https://github.com/kardianos/service)

```
./application run  //or without argument will run
./application install // will install it as a service
./application uninstall // will uninstall the service
./application start // will start the service (you can use the system service manager too)
./application stop // will stopthe service (you can use the system service manager too)

```