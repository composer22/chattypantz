### [Dockerized] (http://www.docker.com) [chattypantz](https://registry.hub.docker.com/u/composer22/chattypantz/)

A docker image for chattypantz. This is created as a single "static" executable using a lightweight image.

To make:

cd docker
./build.sh

Once it completes, you can run the server:
```
docker run --name <containername> -p <hosthttpport>:6660 -d composer22/chattypantz
```

Example below: 6661 and 6061 are host ports e.g. OSX; the others are in the container
We also alter attrs by passing addition params.  Here below, the name of the server is changed from Demo to SanFrancisco, and we also setup the pprof:
```
docker run --name tester -p 6661:6660 -p 6061:6060 -d composer22/chattypantz -N SanFrancisco -L 6060
```
More examples:
```
# Run interactively and watch log output:
 docker run --name tester -ti -p 6660:6660 composer22/chattypantz -N SanFrancisco

# Run as daemon only:
 docker run --name tester -p 6660:6660 c-d omposer22/chattypantz -N SanFrancisco

```
NOTE:  If you are using boot2docker on OSX you also need to map these ports in VirtualBox
if you want to access the server. The easiest way is to launch the OSX app, navigate to the bootdocker-vm, open the settings. From there, select network/adaptor1/advanced/port forwarding. Add an entry such as this:

Example:
```
name: API Server
protocol: TCP
host IP: 127.0.0.1  // this is OSX
host port: 6661   // this is OSX
guest IP: nil
guest port: 6660

save, save
```

Also make sure your /etc/hosts file in OSX has an entry to 127.0.0.1 via localhost if you use this localhost as a name.

#### Options

For additional unix tools in a small image use "FROM busybox" or "FROM phusion/baseimage" instead of "FROM scratch" in Dockerfile.final.
