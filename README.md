# Thesis for University Graduation aa 2020/2021

Built on Hyperledger Fabric and MySQL.

The actual Thesis is written in Italian but this tutorial is in English so anyone can reproduce our work.

### INSTALLATION
> This project is meant to work on Linux/Debian, use that os-family or if you're on a different system DYOR and YMMV, :)

1. Install curl, docker and docker-compose:
```
sudo apt-get update
sudo apt-install curl docker.io -y
curl -L https://github.com/docker/compose/releases/download/1.28.5/docker-compose-`uname -s`-`uname -m` -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
sudo usermod -aG docker ${USER}
newgrp docker
docker run hello-world
```

2. Install and configure Go (version 1.18):
```
wget https://go.dev/dl/go1.18.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.18.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
go version
rm go1.18.linux-amd64.tar.gz
```

3. Install MySQL:
```
sudo apt install mysql-server
sudo mysql_secure_installation
-----no
-----password
-----password
-----no
-----no
-----Y
-----Y
```

4. Configure MySQL:
```
sudo mysql
CREATE USER 'fabric'@'127.0.0.1' IDENTIFIED BY 'password';
CREATE DATABASE tesi;
GRANT ALL PRIVILEGES ON *.* TO 'fabric'@'127.0.0.1' WITH GRANT OPTION;
FLUSH PRIVILEGES;
USE tesi;
CREATE TABLE VeryImportantInfo ( id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY, nome VARCHAR(30) NOT NULL, cognome VARCHAR(30) NOT NULL, saldo VARCHAR(50), data TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP );
exit
```

5. Download Hyperledger Fabric's test net:
```
curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.2.2 1.4.9
cd fabric-samples/fabcar/
```

6. Replace the file "fabcar.go" from the chaincode directory overwriting the same file in fabric-samples/chaincode/fabcar/go directory

7. Start the testnet
```
./startFabric.sh
```

8. Enter the "go" folder and delete all inside
```
cd go
rm -r *
```

9. Replace the client file “fabcar.go” from client directory overwriting the same file in fabric-samples/fabcar/go directory

10. Run the script
```
go mod init fabcar
go mod tidy
go mod vendor
go run fabcar.go
```

11. To close the test net execute the networkDown script in the parent directory
```
cd ..
./networkDown.sh
```
