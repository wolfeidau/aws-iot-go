# iot-aws-go

This is a small project illustrating how to use [golang](https://golang.org/) to provision then connect to the [Amazon Webservices IoT Service](https://aws.amazon.com/iot/).

# prerequisites

Firstly you need an Amazon Webservices (AWS) account to get started, I recommend this guide to get up and running with [Set Up an AWS Account and Create an Administrator User](http://docs.aws.amazon.com/lambda/latest/dg/setting-up.html).

I then recommend following the [aws cli setup](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-set-up.html) to configure your credentials.

Once configured your credentials in the aws cli, and tested them, you are ready to follow the provision section below.

*Note:* Make sure you enable the SES service before attempting to use IoT service, this MUST be activated otherwise you will get an error when you create a thing.

# building

Make sure you have `godep` installed and `$GOPATH/bin` is a in your `$PATH`.

```
go get -u github.com/godep/godep
```

Clone this project into your `$GOPATH`.

```
mkdir -p $GOPATH/src/github.com/wolfeidau
cd $GOPATH/src/github.com/wolfeidau
git clone https://github.com/wolfeidau/aws-iot-go
```

Run make to build and install in `$GOPATH/bin`.

```
make
```

# provision

To provision a new thing based on the preconfigured template.

```
iotprov create --type light --name kitchen_light
```

To create a certificate, policy and link it to the thing you just created, this will save the credentials in `[name].yml` in the current directory. For more information this process [Secure Communication Between a Thing and AWS IoT](https://docs.aws.amazon.com/iot/latest/developerguide/secure-communication.html).

```
iotprov certificate --name kitchen_light
```

To list things.

```
iotprov list --type light
```

Some things to note:

* Create is idempotent as long as you have the matching type
* Currently defaults to us-west-2 as it is closest to Australia
* This will attempt to load the default AWS credentials you loaded when configuring the aws cli or you can export `AWS_PROFILE` to specify one.

# connect

To connect the IoT device using the credentials you just generated.

```
iotdev connect --name kitchen_light --debug
```

# License

iot-aws-go is Copyright (c) 2015 Mark Wolfe @wolfeidau and licensed under the MIT license. All rights not explicitly granted in the MIT license are reserved. See the included LICENSE.md file for more details.
