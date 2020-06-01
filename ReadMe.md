##COMMAND TO RUN
- Install a DOCKER and docker-compose in your System.

Docker-installation
-----------------------------------------------------------------------------------------------------------------------------------
	Sudo Apt-Get Update

	Sudo Apt-Get Install \
	    Apt-Transport-Https \
	    Ca-Certificates \
	    Curl \
	    Gnupg-Agent \
	    Software-Properties-Common

	Curl -Fssl Https://Download.Do

	Cker.Com/Linux/Ubuntu/Gpg | Sudo Apt-Key Add -

	Sudo Apt-Key Fingerprint 0ebfcd88

	Sudo Add-Apt-Repository \
	   "Deb [Arch=Amd64] Https://Download.Docker.Com/Linux/Ubuntu \
	   $(Lsb_Release -Cs) \
	   Stable"


	Sudo Apt-Get Update


	Sudo Apt-Get Install Docker-Ce Docker-Ce-Cli Containerd.Io

-----------------------------------------------------------------------------------------------------------------------------

Docker-compose Installation

	sudo curl -L https://github.com/docker/compose/releases/download/1.18.0/docker-compose-`uname -s`-`uname -m` -o /usr/local/bin/docker-compose

	sudo chmod +x /usr/local/bin/docker-compose
	
------------------------------------

- Run "sudo docker create network oauth"
- Change default user account in dev.yaml file -> environments DEFAULT_USER
- Run "sudo docker-compose -f db.yaml up -d"
- Run "sudo docker-compose -f dev.yaml up"
- Go to browser and type "http://localhost:3000"