#### VWA (Vulnerable Web Application) with GO 
VWA is a web application developed to help the pentester and programmers to learn the vulnerabilities that often occur in web applications which is developed using golang.

#### How To Install VWA

#### Installing docker (for ubuntu 16.04)
First, in order to ensure the downloads are valid, add the GPG key for the official Docker repository to your system:
```bash
$ curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
```
Add the Docker repository to APT sources:
```bash
$ sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
```
Next, update the package database with the Docker packages from the newly added repo:
```bash
$ sudo apt-get update
```
Finally, install Docker:
```bash
$ sudo apt-get install -y docker-ce
```
Docker should now be installed, the daemon started, and the process enabled to start on boot. Check that it's running:
```bash
$ sudo systemctl status docker
```
Output
```bash
docker.service - Docker Application Container Engine
   Loaded: loaded (/lib/systemd/system/docker.service; enabled; vendor preset: enabled)
   Active: active (running) since Thu 2019-01-23 20:28:23 UTC; 35s ago
     Docs: https://docs.docker.com
 Main PID: 13412 (dockerd)
   CGroup: /system.slice/docker.service
           ├─13412 /usr/bin/dockerd -H fd://
           └─13421 docker-containerd --config /var/run/docker/containerd/containerd.toml
```
#### Execute docker command without sudo
Error Output when running docker without sudo:
```
docker: Cannot connect to the Docker daemon. Is the docker daemon running on this host?.
See 'docker run --help'.
```
If you want to avoid typing sudo whenever you run the docker command, add your username to the docker group:
```bash
$ sudo usermod -aG docker ${USER}
```
To apply the new group membership, you can log out of the server and back in, or you can type the following
```bash
$ su - ${USER}
```
You will be prompted to enter your user's password to continue. Afterwards, you can confirm that your user is now added to the docker group by typing:
```
$ id -nG
```
Output
```
yourusername sudo docker
```
#### Installing Docker Compose
We'll check the current release and if necessary, update it in the command below:
```bash
$ sudo curl -L https://github.com/docker/compose/releases/download/1.18.0/docker-compose-`uname -s`-`uname -m` -o /usr/local/bin/docker-compose
```
Next we'll set the permissions:
```bash
$ sudo chmod +x /usr/local/bin/docker-compose
```
Then we'll verify that the installation was successful by checking the version:
```bash
$ docker-compose --version
```

#### Setup Docker to Run VWA
Clone this repository
```bash
$ git clone https://github.com/xCr00tz/vwa.git
```
Change Directory to vwa
```bash
$ cd vwa
```
Run docker compose
```bash
$ docker-compose up
```
#### VWA users

|		email		| password	|
|-------------------|-----------|
| eko@gmail.com		| testing	|
| andi@gmail.com	| testing	|
| attacker@gmail.com| testing	|

Explore the vulnerability. Read the simple pentest report on the folder report/, how can that vulnerability happen & how to mitigate this vulnerability to prevent and patch all vulnerability.

#### To Do

* Reflected & Stored XSS
* IDOR (Insecure Direct Object Reference)
* SQL Injection
* CSRF & Missing CORS Origin

#### Authors
* [Sulhaedir](https://github.com/0c34)
* [Deny](https://github.com/xCr00tz)