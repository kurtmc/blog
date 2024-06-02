---
title: Unifi Autobackup Data Recovery and Restore
published: false
tags: unifi, bash, go
---
# Unifi Autobackup Data Recovery and Restore

[Unifi controller](https://community.ui.com/releases/UniFi-Network-Application-8-0-24/43b24781-aea8-48dc-85b2-3fca42f758c9) is a great piece of software that allows you to easily manage hundreds of WiFi access points, configuring multiple SSIDs, VLAN, etc. The insights dashboard can help you identify clients performance issue and rogue APs. In general I think Unifi is a great tool for home or small business use. So you might have guess, this is a cautionary tale that I hope you don't have to experience yourself since you have a robust backup process that is tested regularly!

Some background: at my work we run Unifi network controller on AWS EC2. There is a redundant backup process, each day we make a full snapshot of the data directory of the Unifi controller application. We do this by first stopping the application, running a simple tar command `tar -czvf $(unifi--data-$(date +%s).tar.gz) config/data`, upload the file to AWS S3 then starting the application back up again. Simple and effective. The restore process has been tested many times and works just as you'd expect. Just for good measure, Unifi provides an incremental backup file that it creates each day too, it creates an autobackup up with a name that looks something like this `autobackup_8.0.24_20240527_1200_1716811200004.unf`, and these are synced to AWS S3 as an insurance policy.

We made a couple mistakes with our S3 backup approach, we put a lifecycle rule on the bucket to delete files after a certain age. The intention was to be able to keep the last `n` backups, but in practice if something goes wrong with the backup process after a few days all the backups will get deleted. We should have implemented a monitor that checks that backups are being produced and alarm early if something is wrong so it can be fixed before any backups get deleted. Hindsight is 20/20.

Some work was being carried out on the EC2 instance that runs Unifi and as part of a normal procedure with these types of third party applications that we run, the instance was terminated and a new one created that would usually look for the latest backup and run the restore procedure automatically. The new instance never restored the data and when it was investigated, we noticed that all the backups were missing.

No worries, the autobackups were still on S3, surely we can just use that right?!

![](https://github.com/kurtmc/blog/raw/master/2024-06/unifi-autobackup-data-recovery-and-restore/images/not_a_valid_backup.png)

```
Error restoring backup
"{filename}" is not a valid backup.
```

*For the purposes of helping out others who might be in this situation, this all relates to version 8.0.24 of Unifi controller and specifically we were using this docker image https://hub.docker.com/r/linuxserver/unifi-controller which is now deprecated and you should be using the following instead https://hub.docker.com/r/linuxserver/unifi-network-application *

**WHAT!**

Ok, there must be something useful in the HTTP response:

```
HTTP/1.1 400 
X-Frame-Options: DENY
Content-Type: application/json;charset=UTF-8
Content-Length: 63
Date: Sun, 02 Jun 2024 03:00:25 GMT
Connection: close
{"meta":{"rc":"error","msg":"api.err.InvalidBackup"},"data":[]}
```

Nope just a 400.

## Can anything be done?

So there is some hope here, I found some things online which can help us inspect the data. First step, there is this repository which provides a bash script to decrypt the backup: https://github.com/zhangyoufu/unifi-backup-decrypt

```
wget https://raw.githubusercontent.com/zhangyoufu/unifi-backup-decrypt/master/decrypt.sh
chmod +x decrypt.sh
./decrypt autobackup_8.0.24_20240527_1200_1716811200004.unf autobackup.zip
```

This produces a zip file, that seems to be completely broken:

```
$ unzip autobackup.zip 
Archive:  autobackup.zip
  End-of-central-directory signature not found.  Either this file is not
  a zipfile, or it constitutes one disk of a multi-part archive.  In the
  latter case the central directory and zipfile comment will be found on
  the last disk(s) of this archive.
unzip:  cannot find zipfile directory in one of autobackup.zip or
        autobackup.zip.zip, and cannot find autobackup.zip.ZIP, period.
```

but, the can be extracted further if you use [7zip](https://www.7-zip.org/download.html), which can be installed on Linux or Mac:

```
# Ubuntu
sudo apt install p7zip-full

# Arch Linux
pacman -S p7zip

# Mac
brew install p7zip
```

Now extract the zip:

```
7z x autobackup.zip
```

This produces `db.gz`, which can further be extacted with `gunzip` to produce a BSON file named `db`

```
gunzip db.gz
```

Now the `db` file can be converted to plain text with MongoDB Database Tools, which can be downloaded from here: https://www.mongodb.com/try/download/database-tools

```
bsondump db > dump.json
```

This may produce a very large file depending on your Unifi deployment size and if you look into the data you see sections and data that seems to be for that section. For for example the lines after `{"__cmd":"select","collection":"devices"}` contain all the mongodb objects in the `devices` collection

I wrote a [small program in go](https://github.com/kurtmc/blog/blob/master/2024-06/unifi-autobackup-data-recovery-and-restore/files/main.go) that can be used to load this data directly into the mongodb database. Which you can do by following these steps:

Copy the `dump.json` file into the container:

```
docker container cp dump.json 5ab135e2d58b:/dump.json
```

Download and run the program:

```
docker exec -it 5ab135e2d58b bash
curl -O https://github.com/kurtmc/blog/raw/master/2024-06/unifi-autobackup-data-recovery-and-restore/files/unifi-restore
./unifi-restore /dump.json
```

Depending on the size of your backup, this can take hours, but once it's complete you can complete, you can restart the Unifi application and you should have access to your data. If the program fails, make sure you read the error, it may be due to the buffer size being too small or the max token size being too small, both of which can be configured using environment variables.

I hope you don't find yourself in this situation where you must rely on the autobackups, but if you do, I hope this helps!
