# iot-aws-go

This is a small project illustrating how to use golang to provision then connect to the [Amazon Webservices IoT Service](https://aws.amazon.com/iot/).

# preperation

To get started with this project you will need to register an [Amazon Webservices](https://aws.amazon.com/) account and download some credentials. I recommend following the [aws cli setup](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-set-up.html).

Once configured your credentials in the aws cli, and tested them, you are ready to follow the provision section below.

Make sure you have `godep` installed and `$GOPATH/bin` is a in your `$PATH`.

```
go get -u github.com/godep/godep
```

# building

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

To create a certificate, policy and link it to the thing you just created, this will save the credentials in `[name].yml` in the current directory.

```
iotprov certificate --name kitchen_light
```

To list things.

```
iotprov list --type light
```

Some things to note:

* Create is idemptotent as long as you have the matching type
* Currently defaults to us-west-2 as it is closest to Australia
* This will attempt to load the default AWS credentials you loaded when configuring the aws cli or you can export `AWS_PROFILE` to specify one.

# connect

To connect the iot device using the credentials you just generated.

```
iotdev connect --name kitchen_light --debug
```

# License

This code is released under the MIT license see the LICENSE.md file for more details.
